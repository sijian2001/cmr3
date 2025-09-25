package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	Email                  string         `gorm:"uniqueIndex;not null" json:"email"`
	Password               string         `gorm:"not null" json:"-"` // "-" タグでJSONレスポンスからパスワードを除外
	Name                   string         `gorm:"not null" json:"name"`
	PasswordResetToken     *string        `gorm:"index" json:"-"` // パスワードリセットトークン
	PasswordResetExpiresAt *time.Time     `json:"-"`              // トークンの有効期限
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}


// TableName はテーブル名を明示的に指定
func (User) TableName() string {
	return "users"
}