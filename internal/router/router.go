package router

import (
	"github.com/RBAC/internal/handler"
	"github.com/gin-gonic/gin"
)

func Start() {
	g := gin.Default()

	// 1. 解决跨域问题 (CORS)
	// 这一步非常重要，否则你的 HTML 文件(如果是本地打开)无法请求后端接口
	g.Use(Cors())

	// 2. 静态文件路由
	// 访问 http://localhost:8080/ 就能看到登录页面
	g.StaticFile("/", "web/login.html")

	// 3. API 路由
	// 前端 axios.post('/login') 会请求到这里
	g.POST("/login", handler.DoLogin)

	// 启动服务 (端口改为 8080 以匹配前端代码)
	g.Run(":8080")
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
