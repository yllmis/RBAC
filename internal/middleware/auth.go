package middleware

import (
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
)

var checkUserPerm = repository.CheckUserPermWithRoute

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
			hasPerm, err := checkUserPerm(userId, requiredPerm, reqMethod, reqPath)
			if err != nil {
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
		cachekey := requiredPerm
		val, err := repository.Rdb.HGet(ctx, redisKey, cachekey).Result()
		if err == nil {
			allowed = (val == "1")
			_ = repository.Rdb.Expire(ctx, redisKey, remaining).Err()
		} else if err == redis.Nil {
			hasPerm, err := checkUserPerm(userId, requiredPerm, reqMethod, reqPath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, tools.ECode{Code: 10001, Message: "internal error"})
				ctx.Abort()
				return
			}
			allowed = hasPerm
			_, err = repository.Rdb.HSet(ctx, redisKey, cachekey, booltoString(allowed)).Result()
			if err == nil {
				if !allowed {
					repository.Rdb.Expire(ctx, redisKey, 30*time.Second)
				} else {
					repository.Rdb.Expire(ctx, redisKey, remaining)
				}
			}
		} else {
			hasPerm, permErr := checkUserPerm(userId, requiredPerm, reqMethod, reqPath)
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

func booltoString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
