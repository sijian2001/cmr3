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

// MockStoreService は StoreServiceInterface のモック
type MockStoreService struct {
	mock.Mock
}

func (m *MockStoreService) CreateStore(store *models.Store) error {
	args := m.Called(store)
	return args.Error(0)
}

func (m *MockStoreService) GetStoreByID(id uint) (*models.Store, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Store), args.Error(1)
}

func (m *MockStoreService) GetStoreByName(name string) (*models.Store, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Store), args.Error(1)
}

func (m *MockStoreService) UpdateStore(store *models.Store) error {
	args := m.Called(store)
	return args.Error(0)
}

func (m *MockStoreService) DeleteStore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStoreService) ListStores(params models.StoreSearchParams) ([]models.Store, int64, error) {
	args := m.Called(params)
	return args.Get(0).([]models.Store), args.Get(1).(int64), args.Error(2)
}

func (m *MockStoreService) UpdateStoreStatus(id uint, status string) (*models.Store, error) {
	args := m.Called(id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Store), args.Error(1)
}

func (m *MockStoreService) GetStoresByStatus(status string) ([]models.Store, error) {
	args := m.Called(status)
	return args.Get(0).([]models.Store), args.Error(1)
}

func (m *MockStoreService) GetActiveStores() ([]models.Store, error) {
	args := m.Called()
	return args.Get(0).([]models.Store), args.Error(1)
}

func (m *MockStoreService) GetInactiveStores() ([]models.Store, error) {
	args := m.Called()
	return args.Get(0).([]models.Store), args.Error(1)
}

func (m *MockStoreService) GetMaintenanceStores() ([]models.Store, error) {
	args := m.Called()
	return args.Get(0).([]models.Store), args.Error(1)
}

func (m *MockStoreService) GetStoreStatusCounts() (*models.StoreStatusCounts, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StoreStatusCounts), args.Error(1)
}

func (m *MockStoreService) GetStoresByManager(managerID uint) ([]models.Store, error) {
	args := m.Called(managerID)
	return args.Get(0).([]models.Store), args.Error(1)
}

func (m *MockStoreService) GetStoresWithoutManager() ([]models.Store, error) {
	args := m.Called()
	return args.Get(0).([]models.Store), args.Error(1)
}

func (m *MockStoreService) AssignManager(storeID, managerID uint) error {
	args := m.Called(storeID, managerID)
	return args.Error(0)
}

func (m *MockStoreService) UnassignManager(storeID uint) error {
	args := m.Called(storeID)
	return args.Error(0)
}

func (m *MockStoreService) BulkUpdateStatus(storeIDs []uint, status string) error {
	args := m.Called(storeIDs, status)
	return args.Error(0)
}

func (m *MockStoreService) IsStoreNameAvailable(name string, excludeID *uint) (bool, error) {
	args := m.Called(name, excludeID)
	return args.Bool(0), args.Error(1)
}

func TestStoreHandler_CreateStore(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な店舗作成",
			requestBody: models.StoreCreateRequest{
				Name:   "テスト店舗",
				Status: "active",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("CreateStore", mock.AnythingOfType("*models.Store")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "無効なリクエスト形式",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockStoreService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "リクエストの形式が正しくありません",
		},
		{
			name: "バリデーションエラー（名前が空）",
			requestBody: models.StoreCreateRequest{
				Name:   "",
				Status: "active",
			},
			mockSetup:      func(m *MockStoreService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "バリデーションエラー",
		},
		{
			name: "店舗名重複エラー",
			requestBody: models.StoreCreateRequest{
				Name:   "重複店舗",
				Status: "active",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("CreateStore", mock.AnythingOfType("*models.Store")).Return(gorm.ErrDuplicatedKey)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "この店舗名は既に存在します",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			var req *http.Request

			if str, ok := tt.requestBody.(string); ok {
				req = httptest.NewRequest("POST", "/stores", bytes.NewBufferString(str))
			} else {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest("POST", "/stores", bytes.NewBuffer(body))
			}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateStore(c)

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

func TestStoreHandler_GetStore(t *testing.T) {
	tests := []struct {
		name           string
		storeID        string
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "正常な店舗取得",
			storeID: "1",
			mockSetup: func(m *MockStoreService) {
				store := &models.Store{
					ID:     1,
					Name:   "テスト店舗",
					Status: "active",
				}
				m.On("GetStoreByID", uint(1)).Return(store, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "無効なID形式",
			storeID:        "invalid",
			mockSetup:      func(m *MockStoreService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "IDの形式が正しくありません",
		},
		{
			name:    "店舗が見つからない",
			storeID: "999",
			mockSetup: func(m *MockStoreService) {
				m.On("GetStoreByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "店舗が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/stores/"+tt.storeID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.storeID)

			err := handler.GetStore(c)

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

func TestStoreHandler_UpdateStore(t *testing.T) {
	tests := []struct {
		name           string
		storeID        string
		requestBody    models.StoreUpdateRequest
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "正常な店舗更新",
			storeID: "1",
			requestBody: models.StoreUpdateRequest{
				Name:   "更新された店舗",
				Status: "inactive",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("UpdateStore", mock.AnythingOfType("*models.Store")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "店舗が見つからない",
			storeID: "999",
			requestBody: models.StoreUpdateRequest{
				Name:   "更新された店舗",
				Status: "active",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("UpdateStore", mock.AnythingOfType("*models.Store")).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "店舗が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/stores/"+tt.storeID, bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.storeID)

			err := handler.UpdateStore(c)

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

func TestStoreHandler_DeleteStore(t *testing.T) {
	tests := []struct {
		name           string
		storeID        string
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "正常な店舗削除",
			storeID: "1",
			mockSetup: func(m *MockStoreService) {
				m.On("DeleteStore", uint(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:    "店舗が見つからない",
			storeID: "999",
			mockSetup: func(m *MockStoreService) {
				m.On("DeleteStore", uint(999)).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "店舗が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("DELETE", "/stores/"+tt.storeID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.storeID)

			err := handler.DeleteStore(c)

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

func TestStoreHandler_ListStores(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "正常な店舗一覧取得",
			queryParams: "page=1&limit=10",
			mockSetup: func(m *MockStoreService) {
				stores := []models.Store{
					{ID: 1, Name: "店舗1", Status: "active"},
					{ID: 2, Name: "店舗2", Status: "inactive"},
				}
				m.On("ListStores", mock.AnythingOfType("models.StoreSearchParams")).Return(stores, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "検索パラメータ付き店舗一覧取得",
			queryParams: "name=テスト&status=active&page=1&limit=5",
			mockSetup: func(m *MockStoreService) {
				stores := []models.Store{
					{ID: 1, Name: "テスト店舗", Status: "active"},
				}
				m.On("ListStores", mock.AnythingOfType("models.StoreSearchParams")).Return(stores, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/stores?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.ListStores(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}

func TestStoreHandler_UpdateStoreStatus(t *testing.T) {
	tests := []struct {
		name           string
		storeID        string
		requestBody    models.StoreStatusUpdateRequest
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "正常なステータス更新",
			storeID: "1",
			requestBody: models.StoreStatusUpdateRequest{
				Status: "inactive",
			},
			mockSetup: func(m *MockStoreService) {
				store := &models.Store{
					ID:     1,
					Name:   "テスト店舗",
					Status: "inactive",
				}
				m.On("UpdateStoreStatus", uint(1), "inactive").Return(store, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "店舗が見つからない",
			storeID: "999",
			requestBody: models.StoreStatusUpdateRequest{
				Status: "inactive",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("UpdateStoreStatus", uint(999), "inactive").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "店舗が見つかりません",
		},
		{
			name:    "無効なステータス遷移",
			storeID: "1",
			requestBody: models.StoreStatusUpdateRequest{
				Status: "active",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("UpdateStoreStatus", uint(1), "active").Return(nil, gorm.ErrInvalidValue)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "無効なステータス遷移です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/stores/"+tt.storeID+"/status", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.storeID)

			err := handler.UpdateStoreStatus(c)

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

func TestStoreHandler_GetStoresByStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "正常なステータス別店舗取得",
			status: "active",
			mockSetup: func(m *MockStoreService) {
				stores := []models.Store{
					{ID: 1, Name: "店舗1", Status: "active"},
					{ID: 2, Name: "店舗2", Status: "active"},
				}
				m.On("GetStoresByStatus", "active").Return(stores, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ステータスパラメータなし",
			status:         "",
			mockSetup:      func(m *MockStoreService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "ステータスパラメータが必要です",
		},
		{
			name:           "無効なステータス",
			status:         "invalid",
			mockSetup:      func(m *MockStoreService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "無効なステータスです",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			url := "/stores/status"
			if tt.status != "" {
				url += "?status=" + tt.status
			}
			req := httptest.NewRequest("GET", url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetStoresByStatus(c)

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

func TestStoreHandler_GetStoreStatusCounts(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なステータス別件数取得",
			mockSetup: func(m *MockStoreService) {
				counts := &models.StoreStatusCounts{
					Active:      5,
					Inactive:    2,
					Maintenance: 1,
					Total:       8,
				}
				m.On("GetStoreStatusCounts").Return(counts, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "サービスエラー",
			mockSetup: func(m *MockStoreService) {
				m.On("GetStoreStatusCounts").Return(nil, fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "ステータス別件数の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			req := httptest.NewRequest("GET", "/stores/status/counts", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetStoreStatusCounts(c)

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

func TestStoreHandler_BulkUpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockStoreService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な一括ステータス更新",
			requestBody: map[string]interface{}{
				"store_ids": []uint{1, 2, 3},
				"status":    "inactive",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("BulkUpdateStatus", []uint{1, 2, 3}, "inactive").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なステータス遷移エラー",
			requestBody: map[string]interface{}{
				"store_ids": []uint{1, 2},
				"status":    "active",
			},
			mockSetup: func(m *MockStoreService) {
				m.On("BulkUpdateStatus", []uint{1, 2}, "active").Return(gorm.ErrInvalidValue)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "無効なステータス遷移が含まれています",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockStoreService{}
			tt.mockSetup(mockService)

			handler := &StoreHandler{
				storeService: mockService,
				validate:     validator.New(),
			}

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/stores/bulk/status", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.BulkUpdateStatus(c)

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