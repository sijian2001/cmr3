package models

import (
	"time"

	"gorm.io/gorm"
)

// Product 製品モデル
type Product struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"not null" validate:"required,min=1,max=100"`
	Description   *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	SKU           string    `json:"sku" gorm:"uniqueIndex;not null" validate:"required,min=1,max=50"`
	Price         float64   `json:"price" gorm:"not null" validate:"required,gt=0,lte=99999999"`
	StockQuantity int       `json:"stock_quantity" gorm:"not null" validate:"gte=0,lte=999999"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ProductCreateRequest 製品作成リクエスト
type ProductCreateRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=100"`
	Description   *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	SKU           string  `json:"sku" validate:"required,min=1,max=50"`
	Price         float64 `json:"price" validate:"required,gt=0,lte=99999999"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0,lte=999999"`
}

// ProductUpdateRequest 製品更新リクエスト
type ProductUpdateRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=100"`
	Description   *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	SKU           string  `json:"sku" validate:"required,min=1,max=50"`
	Price         float64 `json:"price" validate:"required,gt=0,lte=99999999"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0,lte=999999"`
}

// StockAdjustmentRequest 在庫調整リクエスト
type StockAdjustmentRequest struct {
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Operation string `json:"operation" validate:"required,oneof=add subtract"`
	Reason    string `json:"reason,omitempty" validate:"omitempty,max=500"`
}

// ProductSearchParams 製品検索パラメータ
type ProductSearchParams struct {
	Name     string  `query:"name"`
	SKU      string  `query:"sku"`
	MinPrice float64 `query:"min_price"`
	MaxPrice float64 `query:"max_price"`
	Page     int     `query:"page"`
	Limit    int     `query:"limit"`
}

// PaginatedProductResponse ページネーション付き製品レスポンス
type PaginatedProductResponse struct {
	Products   []Product `json:"products"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

// ToProduct ProductCreateRequestをProductに変換
func (req *ProductCreateRequest) ToProduct() Product {
	return Product{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}
}

// ToProduct ProductUpdateRequestをProductに変換
func (req *ProductUpdateRequest) ToProduct() Product {
	return Product{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}
}

// BeforeCreate 作成前のフック
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	// SKUの重複チェック
	var count int64
	if err := tx.Model(&Product{}).Where("sku = ?", p.SKU).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

// BeforeUpdate 更新前のフック
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	// SKUの重複チェック（自分以外）
	var count int64
	if err := tx.Model(&Product{}).Where("sku = ? AND id != ?", p.SKU, p.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

// AdjustStock 在庫数量を調整
func (p *Product) AdjustStock(quantity int, operation string) error {
	switch operation {
	case "add":
		newQuantity := p.StockQuantity + quantity
		if newQuantity > 999999 {
			return gorm.ErrInvalidValue
		}
		p.StockQuantity = newQuantity
	case "subtract":
		newQuantity := p.StockQuantity - quantity
		if newQuantity < 0 {
			p.StockQuantity = 0
		} else {
			p.StockQuantity = newQuantity
		}
	default:
		return gorm.ErrInvalidValue
	}
	return nil
}

// GetStockStatus 在庫ステータスを取得
func (p *Product) GetStockStatus() string {
	if p.StockQuantity == 0 {
		return "在庫切れ"
	} else if p.StockQuantity <= 10 {
		return "在庫少"
	}
	return "在庫十分"
}

// GetStockValue 在庫評価額を計算
func (p *Product) GetStockValue() float64 {
	return p.Price * float64(p.StockQuantity)
}