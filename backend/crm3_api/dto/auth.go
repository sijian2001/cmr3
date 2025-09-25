package dto

// RegisterRequest ユーザー登録リクエスト
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse 認証レスポンス
type AuthResponse struct {
	Token string    `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse ユーザー情報レスポンス
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// PasswordResetRequest パスワードリセットリクエスト
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetConfirmRequest パスワードリセット確認リクエスト
type PasswordResetConfirmRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// MessageResponse メッセージレスポンス
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse エラーレスポンス
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}