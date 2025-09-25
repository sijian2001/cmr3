package utils

import (
	"os"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	// テスト用にJWT_SECRETを設定
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name    string
		userID  uint
		email   string
		wantErr bool
	}{
		{
			name:    "Valid user data",
			userID:  1,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Zero user ID",
			userID:  0,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Empty email",
			userID:  1,
			email:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// トークンが生成されていることを確認
				if len(token) == 0 {
					t.Errorf("GenerateToken() returned empty token")
				}

				// トークンが有効であることを確認
				claims, err := ValidateToken(token)
				if err != nil {
					t.Errorf("Generated token is invalid: %v", err)
				}

				// クレームの内容を確認
				if claims.UserID != tt.userID {
					t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.userID)
				}
				if claims.Email != tt.email {
					t.Errorf("Token Email = %v, want %v", claims.Email, tt.email)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// テスト用にJWT_SECRETを設定
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// 有効なトークンを生成
	userID := uint(1)
	email := "test@example.com"
	validToken, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		wantErr   bool
		wantUser  uint
		wantEmail string
	}{
		{
			name:      "Valid token",
			token:     validToken,
			wantErr:   false,
			wantUser:  userID,
			wantEmail: email,
		},
		{
			name:    "Invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims.UserID != tt.wantUser {
					t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, tt.wantUser)
				}
				if claims.Email != tt.wantEmail {
					t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, tt.wantEmail)
				}

				// 有効期限が正しく設定されていることを確認
				if claims.ExpiresAt == nil {
					t.Errorf("ValidateToken() ExpiresAt is nil")
				} else if time.Now().After(claims.ExpiresAt.Time) {
					t.Errorf("ValidateToken() token is expired")
				}
			}
		})
	}
}

func TestGetJWTSecret(t *testing.T) {
	// 環境変数をクリア
	os.Unsetenv("JWT_SECRET")

	// デフォルト値のテスト
	secret := getJWTSecret()
	if string(secret) != "default-secret-for-development-only" {
		t.Errorf("getJWTSecret() with no env = %v, want default-secret-for-development-only", string(secret))
	}

	// 環境変数設定のテスト
	testSecret := "my-test-secret"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	secret = getJWTSecret()
	if string(secret) != testSecret {
		t.Errorf("getJWTSecret() with env = %v, want %v", string(secret), testSecret)
	}
}