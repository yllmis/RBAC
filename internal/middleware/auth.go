package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RBAC/internal/repository"
	tools "github.com/RBAC/pkg/ecode"
	"github.com/RBAC/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
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
		userId, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
			ctx.Abort()
			return
		}

		// 将 userId 存入上下文，供后续处理使用
		ctx.Set("userId", userId)

		currentPath := ctx.FullPath()

		if currentPath == "" {
			ctx.JSON(http.StatusNotFound, tools.NotFound)
			ctx.Next()
			return
		}
		currentMethod := ctx.Request.Method

		redisField := fmt.Sprintf("%s:%s", currentPath, currentMethod)
		redisKey := fmt.Sprintf("user_prems_%d", userId)

		// 检查用户权限
		exits, err := repository.Rdb.HExists(ctx, redisKey, redisField).Result()

		if err == nil && exits {
			ctx.Next()
			return
		}
		// redis中没有命中，从数据库中查询用户权限

		hasPerm, err := repository.CheckUserPerm(userId, currentPath, currentMethod)
		if err != nil || !hasPerm {
			ctx.JSON(http.StatusForbidden, tools.NoPermission)
			ctx.Abort()
			return
		}
		// 将权限写入redis
		_, err = repository.Rdb.HSet(ctx, redisKey, redisField, 1).Result()
		if err != nil {
			// 写入失败，打印日志，但不影响正常流程
			fmt.Printf("Failed to set user permission in Redis: %v\n", err)
		}

		ctx.Next()

	}
}
