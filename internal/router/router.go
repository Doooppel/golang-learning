package router

import (
	"golang-web-demo/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
// 展示了：
// 1. RESTful API设计
// 2. 路由分组
// 3. 中间件的使用
func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {
	// 创建Gin引擎
	// 如果设置了GIN_MODE=release，会使用生产模式
	r := gin.Default()

	// 添加CORS中间件（跨域支持）
	r.Use(corsMiddleware())

	// API文档路由（简单的HTML页面）
	r.GET("/api/docs", apiDocsHandler)

	// API v1 路由组
	// 展示了路由分组的使用
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)      // POST /api/v1/users - 创建用户
			users.GET("", userHandler.GetAllUsers)      // GET /api/v1/users - 获取所有用户
			users.GET("/:id", userHandler.GetUser)      // GET /api/v1/users/:id - 获取单个用户
			users.PUT("/:id", userHandler.UpdateUser)   // PUT /api/v1/users/:id - 更新用户
			users.DELETE("/:id", userHandler.DeleteUser) // DELETE /api/v1/users/:id - 删除用户
		}

		// 健康检查端点
		api.GET("/health", healthCheckHandler)
	}

	return r
}

// corsMiddleware CORS中间件
// 展示了：
// 1. 中间件的编写
// 2. HTTP头设置
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheckHandler 健康检查处理器
// 展示了简单的响应处理
func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"message": "服务运行正常",
	})
}

// apiDocsHandler API文档处理器
// 展示了HTML响应的返回
func apiDocsHandler(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Golang Web Demo API 文档</title>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
        }
        .endpoint {
            background: white;
            padding: 15px;
            margin: 10px 0;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .method {
            display: inline-block;
            padding: 5px 10px;
            border-radius: 3px;
            color: white;
            font-weight: bold;
            margin-right: 10px;
        }
        .post { background-color: #49cc90; }
        .get { background-color: #61affe; }
        .put { background-color: #fca130; }
        .delete { background-color: #f93e3e; }
        code {
            background-color: #f4f4f4;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: monospace;
        }
        pre {
            background-color: #f4f4f4;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <h1>🚀 Golang Web Demo API 文档</h1>
    
    <div class="endpoint">
        <span class="method get">GET</span>
        <code>/api/v1/health</code>
        <p>健康检查端点</p>
    </div>

    <div class="endpoint">
        <span class="method post">POST</span>
        <code>/api/v1/users</code>
        <p>创建新用户</p>
        <pre>{
  "name": "张三",
  "email": "zhangsan@example.com",
  "age": 25
}</pre>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <code>/api/v1/users</code>
        <p>获取所有用户列表</p>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <code>/api/v1/users/:id</code>
        <p>根据ID获取单个用户</p>
    </div>

    <div class="endpoint">
        <span class="method put">PUT</span>
        <code>/api/v1/users/:id</code>
        <p>更新用户信息（部分更新）</p>
        <pre>{
  "name": "李四",
  "email": "lisi@example.com",
  "age": 30
}</pre>
    </div>

    <div class="endpoint">
        <span class="method delete">DELETE</span>
        <code>/api/v1/users/:id</code>
        <p>删除用户</p>
    </div>

    <h2>📝 使用示例</h2>
    <p>可以使用 curl、Postman 或任何 HTTP 客户端来测试这些 API。</p>
</body>
</html>
`
	c.Data(200, "text/html; charset=utf-8", []byte(html))
}

