package models

import (
	"time"

	"gorm.io/gorm"
)

// Customer 顧客モデル
type Customer struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null;index" validate:"required,min=1,max=100"`
	Email     string         `json:"email" gorm:"not null;uniqueIndex" validate:"required,email,max=255"`
	Phone     *string        `json:"phone,omitempty" gorm:"index" validate:"omitempty,max=50"`
	Address   *string        `json:"address,omitempty" validate:"omitempty,max=500"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName はテーブル名を指定
func (Customer) TableName() string {
	return "customers"
}

// BeforeCreate は作成前のフック
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	// メールアドレスを小文字に正規化
	if c.Email != "" {
		// ここでは簡単な処理のみ実装
		// 実際のプロジェクトではより厳密なバリデーションを行う
	}
	return nil
}

// BeforeUpdate は更新前のフック
func (c *Customer) BeforeUpdate(tx *gorm.DB) error {
	// メールアドレスを小文字に正規化
	if c.Email != "" {
		// ここでは簡単な処理のみ実装
	}
	return nil
}

// CustomerSearchParams 検索パラメータ
type CustomerSearchParams struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Phone    string `form:"phone"`
	Page     int    `form:"page,default=1" validate:"min=1"`
	Limit    int    `form:"limit,default=10" validate:"min=1,max=100"`
	SortBy   string `form:"sort_by,default=id"`
	SortDesc bool   `form:"sort_desc,default=false"`
}

// CustomerCreateRequest 顧客作成リクエスト
type CustomerCreateRequest struct {
	Name    string  `json:"name" validate:"required,min=1,max=100"`
	Email   string  `json:"email" validate:"required,email,max=255"`
	Phone   *string `json:"phone,omitempty" validate:"omitempty,max=50"`
	Address *string `json:"address,omitempty" validate:"omitempty,max=500"`
}

// CustomerUpdateRequest 顧客更新リクエスト
type CustomerUpdateRequest struct {
	Name    string  `json:"name" validate:"required,min=1,max=100"`
	Email   string  `json:"email" validate:"required,email,max=255"`
	Phone   *string `json:"phone,omitempty" validate:"omitempty,max=50"`
	Address *string `json:"address,omitempty" validate:"omitempty,max=500"`
}

// CustomerResponse 顧客レスポンス
type CustomerResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	Address   *string   `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginatedCustomerResponse ページネーション付き顧客リストレスポンス
type PaginatedCustomerResponse struct {
	Customers  []CustomerResponse `json:"customers"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ToResponse モデルをレスポンス用構造体に変換
func (c *Customer) ToResponse() CustomerResponse {
	return CustomerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToCustomer 作成リクエストを顧客モデルに変換
func (req *CustomerCreateRequest) ToCustomer() Customer {
	return Customer{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}
}

// ApplyToCustomer 更新リクエストを顧客モデルに適用
func (req *CustomerUpdateRequest) ApplyToCustomer(customer *Customer) {
	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.Address = req.Address
}