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

var checkUserPerm = repository.CheckUserPerm

func AuthMiddleware(requiredPerm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()
		var allowed bool

		defer func() {
			// 所有退出路径都会记录日志（包括 Abort）
			userId, _ := ctx.Get("userId") // 可能为空（未认证）
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
		val, err := repository.Rdb.HGet(ctx, redisKey, cachekey).Result()
		if err == nil {
			allowed = (val == "1")
			// 同步缓存过期时间为 token 剩余时间
			_ = repository.Rdb.Expire(ctx, redisKey, remaining).Err()
		} else if err == redis.Nil {
			// redis中没有命中，从数据库中查询用户权限

			hasPerm, err := checkUserPerm(userId, requiredPerm)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, tools.ECode{
					Message: "internak error",
				})
				ctx.Abort()
				return
			}
			allowed = hasPerm
			// 将权限写入redis
			_, err = repository.Rdb.HSet(ctx, redisKey, cachekey, booltoString(allowed)).Result()
			if err == nil {
				// 无权限缓存时间要短（比如 30s），避免权限变更后用户长时间无法访问
				if !allowed {
					repository.Rdb.Expire(ctx, redisKey, 30*time.Second)
				} else {
					repository.Rdb.Expire(ctx, redisKey, remaining)
				}
				//“对正向权限和负向权限采用不同 TTL：正向用 Token 剩余时间（长），负向用 30s（短），平衡性能与一致性。”
			}
		} else {
			hasPerm, _ := checkUserPerm(userId, requiredPerm)
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
