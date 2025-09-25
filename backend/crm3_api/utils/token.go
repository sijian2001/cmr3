package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// GenerateResetToken パスワードリセット用トークンを生成
func GenerateResetToken() (string, error) {
	// UUIDベースのトークン生成
	token := uuid.New().String()
	return token, nil
}

// GenerateSecureToken より安全なランダムトークンを生成
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsTokenExpired トークンの有効期限をチェック
func IsTokenExpired(createdAt time.Time, validHours int) bool {
	expirationTime := createdAt.Add(time.Duration(validHours) * time.Hour)
	return time.Now().After(expirationTime)
}