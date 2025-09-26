package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProductCreateRequest_ToProduct(t *testing.T) {
	description := "テスト製品の説明"
	req := ProductCreateRequest{
		Name:          "テスト製品",
		Description:   &description,
		SKU:           "TEST-001",
		Price:         1000.0,
		StockQuantity: 10,
	}

	product := req.ToProduct()

	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.SKU, product.SKU)
	assert.Equal(t, req.Price, product.Price)
	assert.Equal(t, req.StockQuantity, product.StockQuantity)
}

func TestProductUpdateRequest_ToProduct(t *testing.T) {
	description := "更新された製品の説明"
	req := ProductUpdateRequest{
		Name:          "更新された製品",
		Description:   &description,
		SKU:           "TEST-001-UPDATED",
		Price:         1500.0,
		StockQuantity: 15,
	}

	product := req.ToProduct()

	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.SKU, product.SKU)
	assert.Equal(t, req.Price, product.Price)
	assert.Equal(t, req.StockQuantity, product.StockQuantity)
}

func TestProduct_AdjustStock(t *testing.T) {
	tests := []struct {
		name              string
		initialStock      int
		quantity          int
		operation         string
		expectedStock     int
		expectedError     error
	}{
		{
			name:          "正常な在庫追加",
			initialStock:  10,
			quantity:      5,
			operation:     "add",
			expectedStock: 15,
			expectedError: nil,
		},
		{
			name:          "正常な在庫減少",
			initialStock:  10,
			quantity:      3,
			operation:     "subtract",
			expectedStock: 7,
			expectedError: nil,
		},
		{
			name:          "在庫減少で0を下回る場合",
			initialStock:  5,
			quantity:      10,
			operation:     "subtract",
			expectedStock: 0,
			expectedError: nil,
		},
		{
			name:          "在庫追加で上限を超える場合",
			initialStock:  999995,
			quantity:      10,
			operation:     "add",
			expectedStock: 999995, // 変更されない
			expectedError: gorm.ErrInvalidValue,
		},
		{
			name:          "無効な操作",
			initialStock:  10,
			quantity:      5,
			operation:     "invalid",
			expectedStock: 10, // 変更されない
			expectedError: gorm.ErrInvalidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				StockQuantity: tt.initialStock,
			}

			err := product.AdjustStock(tt.quantity, tt.operation)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedStock, product.StockQuantity)
		})
	}
}

func TestProduct_GetStockStatus(t *testing.T) {
	tests := []struct {
		name         string
		stockQty     int
		expectedStatus string
	}{
		{
			name:           "在庫切れ",
			stockQty:       0,
			expectedStatus: "在庫切れ",
		},
		{
			name:           "在庫少（1個）",
			stockQty:       1,
			expectedStatus: "在庫少",
		},
		{
			name:           "在庫少（10個）",
			stockQty:       10,
			expectedStatus: "在庫少",
		},
		{
			name:           "在庫十分（11個）",
			stockQty:       11,
			expectedStatus: "在庫十分",
		},
		{
			name:           "在庫十分（100個）",
			stockQty:       100,
			expectedStatus: "在庫十分",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				StockQuantity: tt.stockQty,
			}

			status := product.GetStockStatus()
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

func TestProduct_GetStockValue(t *testing.T) {
	tests := []struct {
		name          string
		price         float64
		stockQty      int
		expectedValue float64
	}{
		{
			name:          "通常の在庫評価額",
			price:         100.0,
			stockQty:      10,
			expectedValue: 1000.0,
		},
		{
			name:          "在庫ゼロの場合",
			price:         100.0,
			stockQty:      0,
			expectedValue: 0.0,
		},
		{
			name:          "高価格商品",
			price:         9999.99,
			stockQty:      5,
			expectedValue: 49999.95,
		},
		{
			name:          "小数点価格",
			price:         99.99,
			stockQty:      3,
			expectedValue: 299.97,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				Price:         tt.price,
				StockQuantity: tt.stockQty,
			}

			value := product.GetStockValue()
			assert.InDelta(t, tt.expectedValue, value, 0.01) // 小数点誤差を考慮
		})
	}
}

func TestStockAdjustmentRequest_Validation(t *testing.T) {
	// バリデーションはhandlerレベルでテストされるため、
	// ここでは構造体の基本的な作成をテスト
	req := StockAdjustmentRequest{
		Quantity:  10,
		Operation: "add",
		Reason:    "入荷",
	}

	assert.Equal(t, 10, req.Quantity)
	assert.Equal(t, "add", req.Operation)
	assert.Equal(t, "入荷", req.Reason)
}

func TestProductSearchParams_DefaultValues(t *testing.T) {
	// デフォルト値のテストはhandlerレベルで行われるため、
	// ここでは構造体の基本的な作成をテスト
	params := ProductSearchParams{
		Name:     "テスト",
		SKU:      "TEST",
		MinPrice: 100.0,
		MaxPrice: 1000.0,
		Page:     1,
		Limit:    10,
	}

	assert.Equal(t, "テスト", params.Name)
	assert.Equal(t, "TEST", params.SKU)
	assert.Equal(t, 100.0, params.MinPrice)
	assert.Equal(t, 1000.0, params.MaxPrice)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 10, params.Limit)
}

func TestPaginatedProductResponse_Structure(t *testing.T) {
	products := []Product{
		{ID: 1, Name: "製品1", SKU: "TEST-001", Price: 1000.0, StockQuantity: 10},
		{ID: 2, Name: "製品2", SKU: "TEST-002", Price: 2000.0, StockQuantity: 20},
	}

	response := PaginatedProductResponse{
		Products:   products,
		Total:      2,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	assert.Equal(t, 2, len(response.Products))
	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 1, response.TotalPages)
}

// 統合テスト：複数の機能を組み合わせたテスト
func TestProduct_CompleteWorkflow(t *testing.T) {
	// 製品作成リクエストから製品を作成
	description := "テスト用製品"
	createReq := ProductCreateRequest{
		Name:          "テスト製品",
		Description:   &description,
		SKU:           "TEST-WORKFLOW-001",
		Price:         1500.0,
		StockQuantity: 20,
	}

	product := createReq.ToProduct()

	// 初期状態の確認
	assert.Equal(t, "テスト製品", product.Name)
	assert.Equal(t, 20, product.StockQuantity)
	assert.Equal(t, "在庫十分", product.GetStockStatus())
	assert.Equal(t, 30000.0, product.GetStockValue())

	// 在庫調整：追加
	err := product.AdjustStock(5, "add")
	assert.NoError(t, err)
	assert.Equal(t, 25, product.StockQuantity)
	assert.Equal(t, 37500.0, product.GetStockValue())

	// 在庫調整：減少
	err = product.AdjustStock(20, "subtract")
	assert.NoError(t, err)
	assert.Equal(t, 5, product.StockQuantity)
	assert.Equal(t, "在庫少", product.GetStockStatus())
	assert.Equal(t, 7500.0, product.GetStockValue())

	// 在庫調整：在庫切れまで減少
	err = product.AdjustStock(10, "subtract")
	assert.NoError(t, err)
	assert.Equal(t, 0, product.StockQuantity)
	assert.Equal(t, "在庫切れ", product.GetStockStatus())
	assert.Equal(t, 0.0, product.GetStockValue())

	// 更新リクエストで製品情報を更新
	updateReq := ProductUpdateRequest{
		Name:          "更新されたテスト製品",
		Description:   &description,
		SKU:           "TEST-WORKFLOW-001-UPDATED",
		Price:         2000.0,
		StockQuantity: 15,
	}

	updatedProduct := updateReq.ToProduct()
	updatedProduct.ID = product.ID

	assert.Equal(t, "更新されたテスト製品", updatedProduct.Name)
	assert.Equal(t, "TEST-WORKFLOW-001-UPDATED", updatedProduct.SKU)
	assert.Equal(t, 2000.0, updatedProduct.Price)
	assert.Equal(t, 15, updatedProduct.StockQuantity)
	assert.Equal(t, "在庫十分", updatedProduct.GetStockStatus())
	assert.Equal(t, 30000.0, updatedProduct.GetStockValue())
}