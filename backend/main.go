package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/online-editer/backend/handlers"
)

func main() {
	r := gin.Default()

	// CORS設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// コード実行エンドポイント
	r.POST("/execute", handlers.ExecuteCode)

	// 管理用認証API
	r.POST("/admin/auth", handlers.AdminAuth)
	r.GET("/admin/verify", handlers.AdminVerifyToken)

	// 管理用API: runs一覧とファイル取得（認証必須）
	adminAPI := r.Group("/")
	adminAPI.Use(handlers.AdminAuthMiddleware())
	{
		adminAPI.GET("/runs", handlers.ListRuns)
		adminAPI.GET("/runs/:name", handlers.GetRunFile)
		adminAPI.POST("/runs/delete", handlers.DeleteRunFiles)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Backend server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
