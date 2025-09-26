package services

import (
	"crm3_api/models"

	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct 製品を作成
func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.db.Create(product).Error
}

// GetProductByID IDで製品を取得
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := s.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductBySKU SKUで製品を取得
func (s *ProductService) GetProductBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := s.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct 製品を更新
func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.db.Save(product).Error
}

// DeleteProduct 製品を削除
func (s *ProductService) DeleteProduct(id uint) error {
	result := s.db.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// ListProducts 製品一覧を取得（検索とページネーション機能付き）
func (s *ProductService) ListProducts(params models.ProductSearchParams) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{})

	// 検索条件を適用
	if params.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	}
	if params.SKU != "" {
		query = query.Where("sku ILIKE ?", "%"+params.SKU+"%")
	}
	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	// 総数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションを適用して製品を取得
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetStockSummary 在庫サマリーを取得
func (s *ProductService) GetStockSummary() (map[string]interface{}, error) {
	var totalProducts int64
	var lowStockProducts int64
	var outOfStockProducts int64
	var totalStockValue float64

	// 総製品数を取得
	if err := s.db.Model(&models.Product{}).Count(&totalProducts).Error; err != nil {
		return nil, err
	}

	// 在庫少製品数を取得（在庫数が10以下で0より大きい）
	if err := s.db.Model(&models.Product{}).Where("stock_quantity > 0 AND stock_quantity <= 10").Count(&lowStockProducts).Error; err != nil {
		return nil, err
	}

	// 在庫切れ製品数を取得
	if err := s.db.Model(&models.Product{}).Where("stock_quantity = 0").Count(&outOfStockProducts).Error; err != nil {
		return nil, err
	}

	// 総在庫評価額を計算
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}

	for _, product := range products {
		totalStockValue += product.GetStockValue()
	}

	// 在庫状況別の製品数を取得
	stockStatusCounts := make(map[string]int64)
	for _, product := range products {
		status := product.GetStockStatus()
		stockStatusCounts[status]++
	}

	// 価格帯別の製品数を取得
	var priceRanges []struct {
		Range string `json:"range"`
		Count int64  `json:"count"`
	}

	// 価格帯別カウント
	var under1000, from1000to5000, from5000to10000, over10000 int64
	s.db.Model(&models.Product{}).Where("price < 1000").Count(&under1000)
	s.db.Model(&models.Product{}).Where("price >= 1000 AND price < 5000").Count(&from1000to5000)
	s.db.Model(&models.Product{}).Where("price >= 5000 AND price < 10000").Count(&from5000to10000)
	s.db.Model(&models.Product{}).Where("price >= 10000").Count(&over10000)

	priceRanges = append(priceRanges,
		struct {
			Range string `json:"range"`
			Count int64  `json:"count"`
		}{"1000円未満", under1000},
		struct {
			Range string `json:"range"`
			Count int64  `json:"count"`
		}{"1000円-5000円", from1000to5000},
		struct {
			Range string `json:"range"`
			Count int64  `json:"count"`
		}{"5000円-10000円", from5000to10000},
		struct {
			Range string `json:"range"`
			Count int64  `json:"count"`
		}{"10000円以上", over10000},
	)

	summary := map[string]interface{}{
		"total_products":      totalProducts,
		"low_stock_products":  lowStockProducts,
		"out_of_stock_products": outOfStockProducts,
		"total_stock_value":   totalStockValue,
		"stock_status_counts": stockStatusCounts,
		"price_ranges":        priceRanges,
	}

	return summary, nil
}

// GetProductsByStockStatus 在庫ステータス別に製品を取得
func (s *ProductService) GetProductsByStockStatus(status string) ([]models.Product, error) {
	var products []models.Product
	var query *gorm.DB

	switch status {
	case "在庫切れ":
		query = s.db.Where("stock_quantity = 0")
	case "在庫少":
		query = s.db.Where("stock_quantity > 0 AND stock_quantity <= 10")
	case "在庫十分":
		query = s.db.Where("stock_quantity > 10")
	default:
		return nil, gorm.ErrInvalidValue
	}

	if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// GetLowStockProducts 在庫少製品一覧を取得
func (s *ProductService) GetLowStockProducts() ([]models.Product, error) {
	return s.GetProductsByStockStatus("在庫少")
}

// GetOutOfStockProducts 在庫切れ製品一覧を取得
func (s *ProductService) GetOutOfStockProducts() ([]models.Product, error) {
	return s.GetProductsByStockStatus("在庫切れ")
}

// BulkUpdateStock 複数製品の在庫を一括更新
func (s *ProductService) BulkUpdateStock(updates []struct {
	ID       uint `json:"id"`
	Quantity int  `json:"quantity"`
}) error {
	tx := s.db.Begin()

	for _, update := range updates {
		if err := tx.Model(&models.Product{}).Where("id = ?", update.ID).Update("stock_quantity", update.Quantity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}