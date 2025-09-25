package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"crm3_api/config"
	"crm3_api/dto"
	"crm3_api/handlers"
	"crm3_api/middleware"
	"crm3_api/models"
	"crm3_api/utils"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestServer struct {
	Echo *echo.Echo
	DB   *gorm.DB
}

func setupTestServer() (*TestServer, error) {
	// テスト用にJWT_SECRETを設定
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")

	// テスト用のインメモリSQLiteデータベース
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %v", err)
	}

	// テーブル作成
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Echo インスタンス作成
	e := echo.New()

	// ミドルウェア設定（ロギングは省略）
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler(db)

	// パブリック認証API（認証不要）
	api := e.Group("/api")

	// 認証エンドポイント
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/reset-password", authHandler.RequestPasswordReset)
	authGroup.POST("/confirm-reset", authHandler.ConfirmPasswordReset)

	// 認証必須API
	protectedAPI := api.Group("/api")
	protectedAPI.Use(middleware.JWTMiddleware())
	protectedAPI.GET("/user/profile", authHandler.GetProfile)

	return &TestServer{
		Echo: e,
		DB:   db,
	}, nil
}

func (ts *TestServer) Close() {
	os.Unsetenv("JWT_SECRET")
	sqlDB, _ := ts.DB.DB()
	sqlDB.Close()
}

func TestCompleteAuthFlow(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer server.Close()

	// Test data
	userEmail := "integration@example.com"
	userName := "Integration Test User"
	password := "password123"

	t.Run("1. User Registration", func(t *testing.T) {
		registerPayload := dto.RegisterRequest{
			Name:     userName,
			Email:    userEmail,
			Password: password,
		}

		jsonPayload, _ := json.Marshal(registerPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("Registration failed with status %d: %s", rec.Code, rec.Body.String())
		}

		var response dto.AuthResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal registration response: %v", err)
		}

		if response.Token == "" {
			t.Fatal("Registration response missing token")
		}

		if response.User.Email != userEmail {
			t.Errorf("Registration user email = %v, want %v", response.User.Email, userEmail)
		}

		// トークンが有効であることを確認
		claims, err := utils.ValidateToken(response.Token)
		if err != nil {
			t.Fatalf("Registration token is invalid: %v", err)
		}

		if claims.Email != userEmail {
			t.Errorf("Token email = %v, want %v", claims.Email, userEmail)
		}
	})

	t.Run("2. User Login", func(t *testing.T) {
		loginPayload := dto.LoginRequest{
			Email:    userEmail,
			Password: password,
		}

		jsonPayload, _ := json.Marshal(loginPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Login failed with status %d: %s", rec.Code, rec.Body.String())
		}

		var response dto.AuthResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal login response: %v", err)
		}

		if response.Token == "" {
			t.Fatal("Login response missing token")
		}

		// このトークンを後のテストで使用
		os.Setenv("TEST_AUTH_TOKEN", response.Token)
	})

	t.Run("3. Access Protected Endpoint", func(t *testing.T) {
		token := os.Getenv("TEST_AUTH_TOKEN")
		if token == "" {
			t.Fatal("No auth token available from login test")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Profile access failed with status %d: %s", rec.Code, rec.Body.String())
		}

		var response dto.UserResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal profile response: %v", err)
		}

		if response.Email != userEmail {
			t.Errorf("Profile email = %v, want %v", response.Email, userEmail)
		}

		if response.Name != userName {
			t.Errorf("Profile name = %v, want %v", response.Name, userName)
		}
	})

	t.Run("4. Access Protected Endpoint Without Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected unauthorized access, got status %d", rec.Code)
		}
	})

	t.Run("5. Access Protected Endpoint With Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected unauthorized access with invalid token, got status %d", rec.Code)
		}
	})

	t.Run("6. Password Reset Request", func(t *testing.T) {
		resetPayload := dto.PasswordResetRequest{
			Email: userEmail,
		}

		jsonPayload, _ := json.Marshal(resetPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/reset-password", bytes.NewBuffer(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Password reset request failed with status %d: %s", rec.Code, rec.Body.String())
		}

		var response dto.MessageResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal password reset response: %v", err)
		}

		if response.Message == "" {
			t.Fatal("Password reset response missing message")
		}

		// データベースでリセットトークンが設定されていることを確認
		var user models.User
		if err := server.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
			t.Fatalf("Failed to find user: %v", err)
		}

		if user.PasswordResetToken == nil {
			t.Fatal("Password reset token was not set")
		}

		if user.PasswordResetExpiresAt == nil {
			t.Fatal("Password reset expires at was not set")
		}

		// このトークンを次のテストで使用
		os.Setenv("TEST_RESET_TOKEN", *user.PasswordResetToken)
	})

	t.Run("7. Password Reset Confirmation", func(t *testing.T) {
		resetToken := os.Getenv("TEST_RESET_TOKEN")
		if resetToken == "" {
			t.Fatal("No reset token available from password reset request test")
		}

		newPassword := "newpassword456"
		confirmPayload := dto.PasswordResetConfirmRequest{
			Token:    resetToken,
			Password: newPassword,
		}

		jsonPayload, _ := json.Marshal(confirmPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/confirm-reset", bytes.NewBuffer(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Password reset confirmation failed with status %d: %s", rec.Code, rec.Body.String())
		}

		// 新しいパスワードでログインできることを確認
		loginPayload := dto.LoginRequest{
			Email:    userEmail,
			Password: newPassword,
		}

		jsonPayload, _ = json.Marshal(loginPayload)
		req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Login with new password failed with status %d: %s", rec.Code, rec.Body.String())
		}

		// 古いパスワードでログインできないことを確認
		oldLoginPayload := dto.LoginRequest{
			Email:    userEmail,
			Password: password, // 古いパスワード
		}

		jsonPayload, _ = json.Marshal(oldLoginPayload)
		req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()

		server.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected login failure with old password, got status %d", rec.Code)
		}
	})

	// クリーンアップ
	os.Unsetenv("TEST_AUTH_TOKEN")
	os.Unsetenv("TEST_RESET_TOKEN")
}

func TestDuplicateRegistration(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer server.Close()

	userEmail := "duplicate@example.com"
	registerPayload := dto.RegisterRequest{
		Name:     "Duplicate Test User",
		Email:    userEmail,
		Password: "password123",
	}

	// 最初の登録は成功すべき
	jsonPayload, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	server.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("First registration failed with status %d: %s", rec.Code, rec.Body.String())
	}

	// 2回目の登録は失敗すべき（409 Conflict）
	jsonPayload, _ = json.Marshal(registerPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	server.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Expected duplicate registration to fail with 409 Conflict, got status %d", rec.Code)
	}
}