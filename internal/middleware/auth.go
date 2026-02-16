package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RBAC/internal/repository"
	tools "github.com/RBAC/pkg/ecode"
	"github.com/RBAC/pkg/log"
	"github.com/RBAC/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var checkUserPerm = repository.CheckUserPermWithRoute
var permCheckGroup singleflight.Group

func AuthMiddleware(requiredPerm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		var allowed bool

		defer func() {
			userId, _ := ctx.Get("userId")
			log.Logger.Info("rbac_auth_check",
				zap.String("path", ctx.Request.URL.Path),
				zap.String("perm", requiredPerm),
				zap.Any("user_id", userId),
				zap.Bool("allowed", allowed),
				zap.Duration("took", time.Since(start)),
				zap.String("client_ip", ctx.ClientIP()),
			)
		}()

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		userId, expireAt, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		remaining := time.Until(expireAt)
		if remaining <= 0 {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		reqMethod := ctx.Request.Method
		reqPath := ctx.Request.URL.Path

		if repository.Rdb == nil {
			hasPerm, permErr := loadPermWithSingleflight(ctx, userId, requiredPerm, reqMethod, reqPath)
			if permErr != nil {
				ctx.JSON(http.StatusInternalServerError, tools.ECode{Code: 10001, Message: "internal error"})
				ctx.Abort()
				return
			}
			allowed = hasPerm
			if !allowed {
				ctx.JSON(http.StatusForbidden, tools.NoPermission)
				ctx.Abort()
				return
			}
			ctx.Set("userId", userId)
			ctx.Next()
			return
		}

		redisKey := fmt.Sprintf("user_perms_%d", userId)
		cacheKey := requiredPerm
		val, err := repository.Rdb.HGet(ctx, redisKey, cacheKey).Result()
		if err == nil {
			allowed = (val == "1")
			_ = repository.Rdb.Expire(ctx, redisKey, remaining).Err()
		} else if err == redis.Nil {
			hasPerm, permErr := loadPermWithSingleflight(ctx, userId, requiredPerm, reqMethod, reqPath)
			if permErr != nil {
				ctx.JSON(http.StatusInternalServerError, tools.ECode{Code: 10001, Message: "internal error"})
				ctx.Abort()
				return
			}
			allowed = hasPerm

			ttl := remaining
			if !allowed {
				ttl = 30 * time.Second
			}
			if cacheErr := setPermCacheWithTTL(ctx, redisKey, cacheKey, boolToString(allowed), ttl); cacheErr != nil {
				log.Logger.Warn("cache_permission_failed",
					zap.Int64("user_id", userId),
					zap.String("perm", cacheKey),
					zap.Error(cacheErr),
				)
			}
		} else {
			log.Logger.Warn("read_permission_cache_failed",
				zap.Int64("user_id", userId),
				zap.String("perm", cacheKey),
				zap.Error(err),
			)
			hasPerm, permErr := loadPermWithSingleflight(ctx, userId, requiredPerm, reqMethod, reqPath)
			if permErr != nil {
				ctx.JSON(http.StatusInternalServerError, tools.ECode{Code: 10001, Message: "internal error"})
				ctx.Abort()
				return
			}
			allowed = hasPerm
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, tools.NoPermission)
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func loadPermWithSingleflight(ctx context.Context, userID int64, requiredPerm, method, apiPath string) (bool, error) {
	key := fmt.Sprintf("%d|%s|%s|%s", userID, requiredPerm, method, apiPath)
	result, err, _ := permCheckGroup.Do(key, func() (interface{}, error) {
		return checkUserPerm(userID, requiredPerm, method, apiPath)
	})
	if err != nil {
		return false, err
	}

	allowed, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid permission check result type")
	}
	return allowed, nil
}

func setPermCacheWithTTL(ctx context.Context, redisKey, cacheKey, value string, ttl time.Duration) error {
	_, err := repository.Rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, redisKey, cacheKey, value)
		pipe.Expire(ctx, redisKey, ttl)
		return nil
	})
	return err
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
