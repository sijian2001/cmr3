package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"crm3_api/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockProductService は ProductServiceInterface のモック
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) GetProductByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetProductBySKU(sku string) (*models.Product, error) {
	args := m.Called(sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) ListProducts(params models.ProductSearchParams) ([]models.Product, int64, error) {
	args := m.Called(params)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetStockSummary() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockProductService) GetProductsByStockStatus(status string) ([]models.Product, error) {
	args := m.Called(status)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) GetLowStockProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) GetOutOfStockProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) BulkUpdateStock(updates []struct {
	ID       uint `json:"id"`
	Quantity int  `json:"quantity"`
}) error {
	args := m.Called(updates)
	return args.Error(0)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な製品作成",
			requestBody: models.ProductCreateRequest{
				Name:          "テスト製品",
				Description:   nil,
				SKU:           "TEST-001",
				Price:         1000.0,
				StockQuantity: 10,
			},
			mockSetup: func(m *MockProductService) {
				m.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "無効なリクエスト形式",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "リクエストの形式が正しくありません",
		},
		{
			name: "バリデーションエラー（名前が空）",
			requestBody: models.ProductCreateRequest{
				Name:          "",
				SKU:           "TEST-001",
				Price:         1000.0,
				StockQuantity: 10,
			},
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "バリデーションエラー",
		},
		{
			name: "SKU重複エラー",
			requestBody: models.ProductCreateRequest{
				Name:          "テスト製品",
				SKU:           "DUPLICATE-SKU",
				Price:         1000.0,
				StockQuantity: 10,
			},
			mockSetup: func(m *MockProductService) {
				m.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(gorm.ErrDuplicatedKey)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "このSKUは既に存在します",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			var req *http.Request

			if str, ok := tt.requestBody.(string); ok {
				req = httptest.NewRequest("POST", "/products", bytes.NewBufferString(str))
			} else {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
			}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateProduct(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "正常な製品取得",
			productID: "1",
			mockSetup: func(m *MockProductService) {
				product := &models.Product{
					ID:            1,
					Name:          "テスト製品",
					SKU:           "TEST-001",
					Price:         1000.0,
					StockQuantity: 10,
				}
				m.On("GetProductByID", uint(1)).Return(product, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "無効なID形式",
			productID:      "invalid",
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "IDの形式が正しくありません",
		},
		{
			name:      "製品が見つからない",
			productID: "999",
			mockSetup: func(m *MockProductService) {
				m.On("GetProductByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "製品が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/products/"+tt.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			err := handler.GetProduct(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		requestBody    interface{}
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "正常な製品更新",
			productID: "1",
			requestBody: models.ProductUpdateRequest{
				Name:          "更新された製品",
				SKU:           "TEST-001-UPDATED",
				Price:         1500.0,
				StockQuantity: 15,
			},
			mockSetup: func(m *MockProductService) {
				m.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "無効なID形式",
			productID:      "invalid",
			requestBody:    models.ProductUpdateRequest{},
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "IDの形式が正しくありません",
		},
		{
			name:      "製品が見つからない",
			productID: "999",
			requestBody: models.ProductUpdateRequest{
				Name:          "更新された製品",
				SKU:           "TEST-999",
				Price:         1500.0,
				StockQuantity: 15,
			},
			mockSetup: func(m *MockProductService) {
				m.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "製品が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/products/"+tt.productID, bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			err := handler.UpdateProduct(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "正常な製品削除",
			productID: "1",
			mockSetup: func(m *MockProductService) {
				m.On("DeleteProduct", uint(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "無効なID形式",
			productID:      "invalid",
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "IDの形式が正しくありません",
		},
		{
			name:      "製品が見つからない",
			productID: "999",
			mockSetup: func(m *MockProductService) {
				m.On("DeleteProduct", uint(999)).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "製品が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("DELETE", "/products/"+tt.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			err := handler.DeleteProduct(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_ListProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "正常な製品一覧取得",
			queryParams: "page=1&limit=10",
			mockSetup: func(m *MockProductService) {
				products := []models.Product{
					{ID: 1, Name: "製品1", SKU: "TEST-001", Price: 1000.0, StockQuantity: 10},
					{ID: 2, Name: "製品2", SKU: "TEST-002", Price: 2000.0, StockQuantity: 20},
				}
				m.On("ListProducts", mock.AnythingOfType("models.ProductSearchParams")).Return(products, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "検索パラメータ付き製品一覧取得",
			queryParams: "name=テスト&min_price=1000&max_price=5000&page=1&limit=5",
			mockSetup: func(m *MockProductService) {
				products := []models.Product{
					{ID: 1, Name: "テスト製品", SKU: "TEST-001", Price: 1500.0, StockQuantity: 10},
				}
				m.On("ListProducts", mock.AnythingOfType("models.ProductSearchParams")).Return(products, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "デフォルトページネーション",
			queryParams: "",
			mockSetup: func(m *MockProductService) {
				products := []models.Product{}
				m.On("ListProducts", mock.AnythingOfType("models.ProductSearchParams")).Return(products, int64(0), nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/products?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.ListProducts(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_AdjustStock(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		requestBody    models.StockAdjustmentRequest
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "正常な在庫追加",
			productID: "1",
			requestBody: models.StockAdjustmentRequest{
				Quantity:  5,
				Operation: "add",
				Reason:    "入荷",
			},
			mockSetup: func(m *MockProductService) {
				product := &models.Product{
					ID:            1,
					Name:          "テスト製品",
					SKU:           "TEST-001",
					Price:         1000.0,
					StockQuantity: 10,
				}
				m.On("GetProductByID", uint(1)).Return(product, nil)
				m.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "正常な在庫減少",
			productID: "1",
			requestBody: models.StockAdjustmentRequest{
				Quantity:  3,
				Operation: "subtract",
				Reason:    "売上",
			},
			mockSetup: func(m *MockProductService) {
				product := &models.Product{
					ID:            1,
					Name:          "テスト製品",
					SKU:           "TEST-001",
					Price:         1000.0,
					StockQuantity: 10,
				}
				m.On("GetProductByID", uint(1)).Return(product, nil)
				m.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "無効なID形式",
			productID:      "invalid",
			requestBody:    models.StockAdjustmentRequest{},
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "IDの形式が正しくありません",
		},
		{
			name:      "製品が見つからない",
			productID: "999",
			requestBody: models.StockAdjustmentRequest{
				Quantity:  5,
				Operation: "add",
			},
			mockSetup: func(m *MockProductService) {
				m.On("GetProductByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "製品が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/products/"+tt.productID+"/stock", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			err := handler.AdjustStock(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetStockSummary(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な在庫サマリー取得",
			mockSetup: func(m *MockProductService) {
				summary := map[string]interface{}{
					"total_products":      int64(10),
					"low_stock_products":  int64(2),
					"out_of_stock_products": int64(1),
					"total_stock_value":   15000.0,
				}
				m.On("GetStockSummary").Return(summary, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "サービスエラー",
			mockSetup: func(m *MockProductService) {
				m.On("GetStockSummary").Return(map[string]interface{}{}, fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "在庫サマリーの取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/products/stock/summary", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetStockSummary(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_SearchProductsBySKU(t *testing.T) {
	tests := []struct {
		name           string
		skuQuery       string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "正常なSKU検索",
			skuQuery: "TEST-001",
			mockSetup: func(m *MockProductService) {
				product := &models.Product{
					ID:            1,
					Name:          "テスト製品",
					SKU:           "TEST-001",
					Price:         1000.0,
					StockQuantity: 10,
				}
				m.On("GetProductBySKU", "TEST-001").Return(product, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "SKUパラメータなし",
			skuQuery:       "",
			mockSetup:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "SKUパラメータが必要です",
		},
		{
			name:     "SKUが見つからない",
			skuQuery: "NONEXISTENT",
			mockSetup: func(m *MockProductService) {
				m.On("GetProductBySKU", "NONEXISTENT").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "該当するSKUの製品が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{}
			tt.mockSetup(mockService)

			handler := &ProductHandler{
				productService: mockService,
				validate:       validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/products/search/sku?sku="+tt.skuQuery, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.SearchProductsBySKU(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}