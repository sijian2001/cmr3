package main

import (
	"log"
	"net/http"

	"crm3_api/config"
	"crm3_api/database"
	"crm3_api/handlers"
	"crm3_api/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var db *gorm.DB

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "CRM3 API Server is running",
		"status":  "ok",
	})
}


func main() {
	// データベース接続
	var err error
	db, err = config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// マイグレーション実行
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Echo インスタンス作成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.RequestResponseLogging())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler(db)

	// パブリックルート
	e.GET("/", hello)

	// パブリック認証API（認証不要）
	api := e.Group("/api")

	// 認証エンドポイント用のレートリミット
	authGroup := api.Group("/auth")
	authGroup.Use(middleware.AuthRateLimiter())
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// パスワードリセット用のより厳しいレートリミット
	resetGroup := api.Group("/auth")
	resetGroup.Use(middleware.PasswordResetRateLimiter())
	resetGroup.POST("/reset-password", authHandler.RequestPasswordReset)
	resetGroup.POST("/confirm-reset", authHandler.ConfirmPasswordReset)

	// 認証必須API
	protectedAPI := e.Group("/api")
	protectedAPI.Use(middleware.JWTMiddleware())
	protectedAPI.GET("/user/profile", authHandler.GetProfile)

	// サーバー起動
	log.Println("Starting server on :1323")
	e.Start(":1323")
}