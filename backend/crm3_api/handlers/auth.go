package handlers

import (
	"log"
	"net/http"
	"time"

	"crm3_api/dto"
	"crm3_api/models"
	"crm3_api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Validator: validator.New(),
	}
}

// Register ユーザー登録
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// メールアドレスの重複チェック
	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "email_exists",
			Message: "このメールアドレスは既に登録されています",
		})
	}

	// パスワードハッシュ化
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Password hashing failed for email %s: %v", req.Email, err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "パスワードの処理に失敗しました",
		})
	}

	// トランザクション内でユーザー作成
	var user models.User
	var token string
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// ユーザー作成
		user = models.User{
			Email:    req.Email,
			Password: hashedPassword,
			Name:     req.Name,
		}

		if err := tx.Create(&user).Error; err != nil {
			log.Printf("User creation failed for email %s: %v", req.Email, err)
			return err
		}

		// JWTトークン生成
		token, err = utils.GenerateToken(user.ID, user.Email)
		if err != nil {
			log.Printf("JWT token generation failed for user ID %d: %v", user.ID, err)
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "ユーザーの作成に失敗しました",
		})
	}

	log.Printf("User registered successfully: %s (ID: %d)", user.Email, user.ID)
	return c.JSON(http.StatusCreated, dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

// Login ログイン
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// ユーザー検索
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		log.Printf("Login failed for email %s: user not found", req.Email)
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "メールアドレスまたはパスワードが正しくありません",
		})
	}

	// パスワード確認
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		log.Printf("Login failed for email %s: invalid password", req.Email)
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "メールアドレスまたはパスワードが正しくありません",
		})
	}

	// JWTトークン生成
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "token_error",
			Message: "認証トークンの生成に失敗しました",
		})
	}

	log.Printf("User logged in successfully: %s (ID: %d)", user.Email, user.ID)
	return c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

// GetProfile ユーザープロフィール取得（認証必須）
func (h *AuthHandler) GetProfile(c echo.Context) error {
	// JWT認証ミドルウェアで設定されたユーザー情報を取得
	userID := c.Get("user_id").(uint)

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "user_not_found",
			Message: "ユーザーが見つかりません",
		})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// RequestPasswordReset パスワードリセットリクエスト
func (h *AuthHandler) RequestPasswordReset(c echo.Context) error {
	var req dto.PasswordResetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// ユーザー検索
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// セキュリティ上、ユーザーが存在しない場合でも成功レスポンスを返す
		return c.JSON(http.StatusOK, dto.MessageResponse{
			Message: "パスワードリセットメールを送信しました。メールボックスをご確認ください。",
		})
	}

	// リセットトークン生成
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "token_error",
			Message: "リセットトークンの生成に失敗しました",
		})
	}

	// トークン有効期限設定（1時間）
	expiresAt := time.Now().Add(1 * time.Hour)

	// ユーザーレコードを更新
	if err := h.DB.Model(&user).Updates(map[string]interface{}{
		"password_reset_token":     resetToken,
		"password_reset_expires_at": expiresAt,
	}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "トークンの保存に失敗しました",
		})
	}

	// パスワードリセットメール送信
	if err := utils.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "email_error",
			Message: "メール送信に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "パスワードリセットメールを送信しました。メールボックスをご確認ください。",
	})
}

// ConfirmPasswordReset パスワードリセット確認・実行
func (h *AuthHandler) ConfirmPasswordReset(c echo.Context) error {
	var req dto.PasswordResetConfirmRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// トークンでユーザー検索
	var user models.User
	if err := h.DB.Where("password_reset_token = ?", req.Token).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_token",
			Message: "無効なリセットトークンです",
		})
	}

	// トークン有効期限チェック
	if user.PasswordResetExpiresAt == nil || time.Now().After(*user.PasswordResetExpiresAt) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "token_expired",
			Message: "リセットトークンの有効期限が切れています",
		})
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "パスワードの処理に失敗しました",
		})
	}

	// トランザクション内でパスワード更新・トークンクリア
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&user).Updates(map[string]interface{}{
			"password":                   hashedPassword,
			"password_reset_token":       nil,
			"password_reset_expires_at": nil,
		}).Error
	})

	if err != nil {
		log.Printf("Password update failed for user ID %d: %v", user.ID, err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "パスワードの更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "パスワードが正常に更新されました",
	})
}