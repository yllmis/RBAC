package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RBAC/internal/repository"
	tools "github.com/RBAC/pkg/ecode"
	"github.com/RBAC/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(requiredPerm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		// 兼容 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		// 解析 token，获取userId
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

		redisKey := fmt.Sprintf("user_perms_%d", userId)
		cachekey := requiredPerm
		// 检查用户权限
		var allowed bool
		val, err := repository.Rdb.HGet(ctx, redisKey, cachekey).Result()
		if err == nil {
			allowed = (val == "1")
			// 同步缓存过期时间为 token 剩余时间
			_ = repository.Rdb.Expire(ctx, redisKey, remaining).Err()
		} else if err == redis.Nil {
			// redis中没有命中，从数据库中查询用户权限

			hasPerm, err := repository.CheckUserPerm(userId, requiredPerm)
			if err != nil || !hasPerm {
				ctx.JSON(http.StatusInternalServerError, tools.NoPermission)
				ctx.Abort()
				return
			}
			allowed = hasPerm
			// 将权限写入redis
			_, err = repository.Rdb.HSet(ctx, redisKey, cachekey, booltoString(allowed)).Result()
			if err == nil {
				_, _ = repository.Rdb.Expire(ctx, redisKey, remaining).Result()
			}
		} else {
			hasPerm, _ := repository.CheckUserPerm(userId, requiredPerm)
			allowed = hasPerm
		}
		if !allowed {
			ctx.JSON(http.StatusForbidden, tools.NoPermission)
			ctx.Abort()
			return
		}

		// 将 userId 存入上下文，供后续处理使用
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
