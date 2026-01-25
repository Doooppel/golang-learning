package main

import (
	"log"
	"os"

	"golang-web-demo/internal/database"
	"golang-web-demo/internal/handlers"
	"golang-web-demo/internal/router"
)

func main() {
	// 基础语法示例：变量声明和初始化
	// Go语言基础：变量声明方式
	var port string = "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// 基础语法：创建内存数据库实例
	// 这是一个单例模式的应用
	db := database.NewInMemoryDB()

	// 基础语法：初始化处理器
	userHandler := handlers.NewUserHandler(db)

	// 基础语法：创建Gin路由引擎
	// Gin是Go语言最流行的Web框架之一
	r := router.SetupRouter(userHandler)

	// 基础语法：启动HTTP服务器
	// 这里展示了Go语言的错误处理模式
	log.Printf("🚀 服务器启动在端口 %s", port)
	log.Printf("📖 API文档: http://localhost:%s/api/docs", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
