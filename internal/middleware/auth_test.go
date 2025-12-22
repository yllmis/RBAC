package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RBAC/internal/repository"
	"github.com/RBAC/pkg/log"
	"github.com/RBAC/pkg/utils"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func setupRouter(requiredPerm string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ping", AuthMiddleware(requiredPerm), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func makeRequest(r *gin.Engine, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func withStubRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	repository.Rdb = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return mr
}

func TestAuth_CacheHit(t *testing.T) {
	log.Init()
	mr := withStubRedis(t)
	defer mr.Close()

	userID := int64(1)
	token, err := utils.GenerateToken(userID)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	key := fmt.Sprintf("user_perms_%d", userID)
	mr.HSet(key, "user:list", "1")

	r := setupRouter("user:list")
	w := makeRequest(r, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAuth_CacheMiss_DBAllow(t *testing.T) {
	log.Init()
	mr := withStubRedis(t)
	defer mr.Close()

	original := checkUserPerm
	checkUserPerm = func(int64, string) (bool, error) { return true, nil }
	t.Cleanup(func() { checkUserPerm = original })

	userID := int64(2)
	token, err := utils.GenerateToken(userID)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := setupRouter("user:list")
	w := makeRequest(r, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	key := fmt.Sprintf("user_perms_%d", userID)
	val := mr.HGet(key, "user:list")
	if val != "1" {
		t.Fatalf("expected cached value '1', got %q", val)
	}
	if ttl := mr.TTL(key); ttl <= 0 {
		t.Fatalf("expected TTL to be set, got %v", ttl)
	}
}

func TestAuth_PermissionDenied(t *testing.T) {
	log.Init()
	mr := withStubRedis(t)
	defer mr.Close()

	original := checkUserPerm
	checkUserPerm = func(int64, string) (bool, error) { return false, nil }
	t.Cleanup(func() { checkUserPerm = original })

	token, err := utils.GenerateToken(3)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := setupRouter("user:list")
	w := makeRequest(r, token)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	log.Init()
	mr := withStubRedis(t)
	defer mr.Close()

	r := setupRouter("user:list")
	w := makeRequest(r, "bad-token")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
