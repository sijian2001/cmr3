package main

import (
	"log"
	"net/http"

	"crm3_api/config"
	"crm3_api/database"
	"crm3_api/handlers"
	"crm3_api/middleware"
	"crm3_api/services"

	"github.com/go-playground/validator/v10"
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
	validate := validator.New()
	authHandler := handlers.NewAuthHandler(db)
	customerHandler := handlers.NewCustomerHandler(db)

	// サービス初期化
	productService := services.NewProductService(db)
	productHandler := handlers.NewProductHandler(productService, validate)

	storeService := services.NewStoreService(db)
	storeHandler := handlers.NewStoreHandler(storeService, validate)

	staffService := services.NewStaffService(db)
	staffHandler := handlers.NewStaffHandler(staffService)

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

	// 顧客管理API
	customerAPI := protectedAPI.Group("/customers")
	customerAPI.GET("", customerHandler.GetCustomers)         // 顧客一覧取得
	customerAPI.POST("", customerHandler.CreateCustomer)      // 顧客作成
	customerAPI.GET("/:id", customerHandler.GetCustomer)      // 顧客詳細取得
	customerAPI.PUT("/:id", customerHandler.UpdateCustomer)   // 顧客更新
	customerAPI.DELETE("/:id", customerHandler.DeleteCustomer) // 顧客削除

	// 製品管理API
	productAPI := protectedAPI.Group("/products")
	productAPI.GET("", productHandler.ListProducts)           // 製品一覧取得（検索・ページネーション付き）
	productAPI.POST("", productHandler.CreateProduct)         // 製品作成
	productAPI.GET("/:id", productHandler.GetProduct)         // 製品詳細取得
	productAPI.PUT("/:id", productHandler.UpdateProduct)      // 製品更新
	productAPI.DELETE("/:id", productHandler.DeleteProduct)   // 製品削除
	productAPI.POST("/:id/stock", productHandler.AdjustStock) // 在庫調整
	productAPI.GET("/stock/summary", productHandler.GetStockSummary) // 在庫サマリー取得
	productAPI.GET("/search/sku", productHandler.SearchProductsBySKU) // SKU検索

	// 店舗管理API
	storeAPI := protectedAPI.Group("/stores")
	storeAPI.GET("", storeHandler.ListStores)                    // 店舗一覧取得（検索・ページネーション付き）
	storeAPI.POST("", storeHandler.CreateStore)                  // 店舗作成
	storeAPI.GET("/:id", storeHandler.GetStore)                  // 店舗詳細取得
	storeAPI.PUT("/:id", storeHandler.UpdateStore)               // 店舗更新
	storeAPI.DELETE("/:id", storeHandler.DeleteStore)            // 店舗削除
	storeAPI.POST("/:id/status", storeHandler.UpdateStoreStatus) // ステータス更新
	storeAPI.GET("/status", storeHandler.GetStoresByStatus)      // ステータス別店舗一覧
	storeAPI.GET("/status/counts", storeHandler.GetStoreStatusCounts) // ステータス別件数
	storeAPI.POST("/:id/manager", storeHandler.AssignManager)    // 管理者割り当て
	storeAPI.DELETE("/:id/manager", storeHandler.UnassignManager) // 管理者割り当て解除
	storeAPI.GET("/manager", storeHandler.GetStoresByManager)    // 管理者別店舗一覧
	storeAPI.POST("/bulk/status", storeHandler.BulkUpdateStatus) // ステータス一括更新
	storeAPI.GET("/check/name", storeHandler.CheckStoreName)     // 店舗名チェック

	// スタッフ管理API
	staffAPI := protectedAPI.Group("/staff")
	staffAPI.GET("", staffHandler.GetAllStaff)                      // スタッフ一覧取得（検索・ページネーション付き）
	staffAPI.POST("", staffHandler.CreateStaff)                     // スタッフ作成
	staffAPI.GET("/:id", staffHandler.GetStaffByID)                 // スタッフ詳細取得
	staffAPI.PUT("/:id", staffHandler.UpdateStaff)                  // スタッフ更新
	staffAPI.DELETE("/:id", staffHandler.DeleteStaff)               // スタッフ削除
	staffAPI.POST("/:id/status", staffHandler.UpdateStaffStatus)    // ステータス更新
	staffAPI.POST("/:id/store", staffHandler.AssignToStore)         // 店舗割り当て
	staffAPI.DELETE("/:id/store", staffHandler.UnassignFromStore)   // 店舗割り当て解除
	staffAPI.GET("/status-counts", staffHandler.GetStaffStatusCounts) // ステータス別件数
	staffAPI.GET("/unassigned", staffHandler.GetUnassignedStaff)    // 未割り当てスタッフ一覧
	staffAPI.POST("/bulk-update-status", staffHandler.BulkUpdateStatus) // ステータス一括更新
	staffAPI.POST("/bulk-assign-store", staffHandler.BulkAssignToStore) // 店舗一括割り当て
	staffAPI.GET("/search", staffHandler.SearchStaffByName)         // 名前検索

	// 店舗別スタッフAPI
	protectedAPI.GET("/stores/:store_id/staff", staffHandler.GetStaffByStore) // 店舗別スタッフ一覧

	// サーバー起動
	log.Println("Starting server on :1323")
	e.Start(":1323")
}