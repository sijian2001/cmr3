package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"crm3_api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStaffService is a mock implementation of StaffService
type MockStaffService struct {
	mock.Mock
}

func (m *MockStaffService) GetAllStaff(params models.StaffSearchParams) (*models.PaginatedStaffResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedStaffResponse), args.Error(1)
}

func (m *MockStaffService) GetStaffByID(id uint) (*models.Staff, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockStaffService) CreateStaff(request models.StaffCreateRequest) (*models.Staff, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockStaffService) UpdateStaff(id uint, request models.StaffUpdateRequest) (*models.Staff, error) {
	args := m.Called(id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockStaffService) DeleteStaff(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStaffService) UpdateStaffStatus(id uint, status string) (*models.Staff, error) {
	args := m.Called(id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockStaffService) AssignToStore(id uint, storeID uint) (*models.Staff, error) {
	args := m.Called(id, storeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockStaffService) UnassignFromStore(id uint) (*models.Staff, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockStaffService) GetStaffStatusCounts() (*models.StaffStatusCounts, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StaffStatusCounts), args.Error(1)
}

func (m *MockStaffService) GetStaffByStore(storeID uint) ([]models.Staff, error) {
	args := m.Called(storeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Staff), args.Error(1)
}

func (m *MockStaffService) GetUnassignedStaff() ([]models.Staff, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Staff), args.Error(1)
}

func (m *MockStaffService) BulkUpdateStatus(ids []uint, status string) error {
	args := m.Called(ids, status)
	return args.Error(0)
}

func (m *MockStaffService) BulkAssignToStore(ids []uint, storeID uint) error {
	args := m.Called(ids, storeID)
	return args.Error(0)
}

func (m *MockStaffService) IsEmailAvailable(email string, excludeID *uint) (bool, error) {
	args := m.Called(email, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStaffService) SearchStaffByName(name string) ([]models.Staff, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Staff), args.Error(1)
}

func TestStaffHandler_GetAllStaff(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	tests := []struct {
		name           string
		queryParams    string
		mockReturn     *models.PaginatedStaffResponse
		mockError      error
		expectedStatus int
	}{
		{
			name:        "成功：スタッフ一覧取得",
			queryParams: "?page=1&limit=10",
			mockReturn: &models.PaginatedStaffResponse{
				Staff: []models.Staff{
					{ID: 1, Name: "田中太郎", Email: "tanaka@example.com", Position: "店長", Status: "active"},
					{ID: 2, Name: "佐藤花子", Email: "sato@example.com", Position: "正社員", Status: "inactive"},
				},
				Total:      2,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：データベースエラー",
			queryParams:    "",
			mockReturn:     nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/staff"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Clear previous expectations
			mockService.ExpectedCalls = nil
			mockService.On("GetAllStaff", mock.AnythingOfType("models.StaffSearchParams")).Return(tt.mockReturn, tt.mockError)

			err := handler.GetAllStaff(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.mockReturn != nil {
				var response models.PaginatedStaffResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockReturn.Staff), len(response.Staff))
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestStaffHandler_GetStaffByID(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	tests := []struct {
		name           string
		staffID        string
		mockReturn     *models.Staff
		mockError      error
		expectedStatus int
	}{
		{
			name:    "成功：スタッフ詳細取得",
			staffID: "1",
			mockReturn: &models.Staff{
				ID:       1,
				Name:     "田中太郎",
				Email:    "tanaka@example.com",
				Position: "店長",
				Status:   "active",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：スタッフが見つからない",
			staffID:        "999",
			mockReturn:     nil,
			mockError:      errors.New("staff not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "エラー：無効なID",
			staffID:        "invalid",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/staff/"+tt.staffID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.staffID)

			if tt.staffID != "invalid" {
				id, _ := strconv.ParseUint(tt.staffID, 10, 32)
				mockService.On("GetStaffByID", uint(id)).Return(tt.mockReturn, tt.mockError)
			}

			err := handler.GetStaffByID(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.staffID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestStaffHandler_CreateStaff(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	phone := "090-1234-5678"
	storeID := uint(1)

	validRequest := models.StaffCreateRequest{
		Name:     "田中太郎",
		Email:    "tanaka@example.com",
		Phone:    &phone,
		Position: "店長",
		StoreID:  &storeID,
		HireDate: "2023-01-01",
		Status:   "active",
	}

	createdStaff := &models.Staff{
		ID:       1,
		Name:     validRequest.Name,
		Email:    validRequest.Email,
		Phone:    validRequest.Phone,
		Position: validRequest.Position,
		StoreID:  validRequest.StoreID,
		HireDate: validRequest.HireDate,
		Status:   validRequest.Status,
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		mockReturn     *models.Staff
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：スタッフ作成",
			requestBody:    validRequest,
			mockReturn:     createdStaff,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "エラー：メールアドレス重複",
			requestBody:    validRequest,
			mockReturn:     nil,
			mockError:      errors.New("email already exists"),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "エラー：店舗が見つからない",
			requestBody:    validRequest,
			mockReturn:     nil,
			mockError:      errors.New("store not found"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "エラー：無効なリクエスト",
			requestBody:    "invalid json",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			// Clear previous expectations
			mockService.ExpectedCalls = nil

			if requestBody, ok := tt.requestBody.(models.StaffCreateRequest); ok {
				json.NewEncoder(&body).Encode(requestBody)
				mockService.On("CreateStaff", requestBody).Return(tt.mockReturn, tt.mockError)
			} else {
				body.WriteString(tt.requestBody.(string))
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/staff", &body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateStaff(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if _, ok := tt.requestBody.(models.StaffCreateRequest); ok {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestStaffHandler_UpdateStaff(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	phone := "090-5678-1234"
	storeID := uint(2)

	validRequest := models.StaffUpdateRequest{
		Name:     "佐藤花子",
		Email:    "sato@example.com",
		Phone:    &phone,
		Position: "副店長",
		StoreID:  &storeID,
		HireDate: "2023-02-01",
		Status:   "inactive",
	}

	updatedStaff := &models.Staff{
		ID:       1,
		Name:     validRequest.Name,
		Email:    validRequest.Email,
		Phone:    validRequest.Phone,
		Position: validRequest.Position,
		StoreID:  validRequest.StoreID,
		HireDate: validRequest.HireDate,
		Status:   validRequest.Status,
	}

	tests := []struct {
		name           string
		staffID        string
		requestBody    interface{}
		mockReturn     *models.Staff
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：スタッフ更新",
			staffID:        "1",
			requestBody:    validRequest,
			mockReturn:     updatedStaff,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：スタッフが見つからない",
			staffID:        "999",
			requestBody:    validRequest,
			mockReturn:     nil,
			mockError:      errors.New("staff not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "エラー：無効なID",
			staffID:        "invalid",
			requestBody:    validRequest,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/staff/"+tt.staffID, &body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.staffID)

			if tt.staffID != "invalid" {
				id, _ := strconv.ParseUint(tt.staffID, 10, 32)
				mockService.On("UpdateStaff", uint(id), validRequest).Return(tt.mockReturn, tt.mockError)
			}

			err := handler.UpdateStaff(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.staffID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestStaffHandler_DeleteStaff(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	tests := []struct {
		name           string
		staffID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：スタッフ削除",
			staffID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：スタッフが見つからない",
			staffID:        "999",
			mockError:      errors.New("staff not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "エラー：無効なID",
			staffID:        "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/staff/"+tt.staffID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.staffID)

			if tt.staffID != "invalid" {
				id, _ := strconv.ParseUint(tt.staffID, 10, 32)
				mockService.On("DeleteStaff", uint(id)).Return(tt.mockError)
			}

			err := handler.DeleteStaff(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.staffID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestStaffHandler_UpdateStaffStatus(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	updatedStaff := &models.Staff{
		ID:     1,
		Name:   "田中太郎",
		Status: "inactive",
	}

	tests := []struct {
		name           string
		staffID        string
		requestBody    map[string]string
		mockReturn     *models.Staff
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：ステータス更新",
			staffID:        "1",
			requestBody:    map[string]string{"status": "inactive"},
			mockReturn:     updatedStaff,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：スタッフが見つからない",
			staffID:        "999",
			requestBody:    map[string]string{"status": "inactive"},
			mockReturn:     nil,
			mockError:      errors.New("staff not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "エラー：無効なステータス",
			staffID:        "1",
			requestBody:    map[string]string{"status": "invalid_status"},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/staff/"+tt.staffID+"/status", &body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.staffID)

			if tt.staffID != "invalid" && tt.requestBody["status"] != "invalid_status" {
				id, _ := strconv.ParseUint(tt.staffID, 10, 32)
				mockService.On("UpdateStaffStatus", uint(id), tt.requestBody["status"]).Return(tt.mockReturn, tt.mockError)
			}

			err := handler.UpdateStaffStatus(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.staffID != "invalid" && tt.requestBody["status"] != "invalid_status" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestStaffHandler_AssignToStore(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	assignedStaff := &models.Staff{
		ID:      1,
		Name:    "田中太郎",
		StoreID: func() *uint { id := uint(1); return &id }(),
	}

	tests := []struct {
		name           string
		staffID        string
		requestBody    map[string]uint
		mockReturn     *models.Staff
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：店舗割り当て",
			staffID:        "1",
			requestBody:    map[string]uint{"store_id": 1},
			mockReturn:     assignedStaff,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：スタッフが見つからない",
			staffID:        "999",
			requestBody:    map[string]uint{"store_id": 1},
			mockReturn:     nil,
			mockError:      errors.New("staff not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "エラー：店舗が見つからない",
			staffID:        "1",
			requestBody:    map[string]uint{"store_id": 999},
			mockReturn:     nil,
			mockError:      errors.New("store not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/staff/"+tt.staffID+"/store", &body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.staffID)

			if tt.staffID != "invalid" {
				id, _ := strconv.ParseUint(tt.staffID, 10, 32)
				mockService.On("AssignToStore", uint(id), tt.requestBody["store_id"]).Return(tt.mockReturn, tt.mockError)
			}

			err := handler.AssignToStore(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.staffID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestStaffHandler_GetStaffStatusCounts(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	mockCounts := &models.StaffStatusCounts{
		Active:   5,
		Inactive: 2,
		OnLeave:  1,
		Total:    8,
	}

	tests := []struct {
		name           string
		mockReturn     *models.StaffStatusCounts
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：ステータス別件数取得",
			mockReturn:     mockCounts,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：データベースエラー",
			mockReturn:     nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/staff/status-counts", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Clear previous expectations
			mockService.ExpectedCalls = nil
			mockService.On("GetStaffStatusCounts").Return(tt.mockReturn, tt.mockError)

			err := handler.GetStaffStatusCounts(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.mockReturn != nil {
				var response models.StaffStatusCounts
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn.Active, response.Active)
				assert.Equal(t, tt.mockReturn.Total, response.Total)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestStaffHandler_BulkUpdateStatus(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	validRequest := map[string]interface{}{
		"ids":    []uint{1, 2, 3},
		"status": "inactive",
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：ステータス一括更新",
			requestBody:    validRequest,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：データベースエラー",
			requestBody:    validRequest,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "エラー：無効なリクエスト",
			requestBody:    map[string]interface{}{"invalid": "request"},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/staff/bulk-update-status", &body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if request, ok := tt.requestBody.(map[string]interface{}); ok {
				if ids, exists := request["ids"]; exists {
					if idsSlice, ok := ids.([]uint); ok {
						mockService.On("BulkUpdateStatus", idsSlice, request["status"].(string)).Return(tt.mockError)
					}
				}
			}

			err := handler.BulkUpdateStatus(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Only assert expectations if we set up the mock
			if _, ok := tt.requestBody.(map[string]interface{}); ok {
				if request, ok := tt.requestBody.(map[string]interface{}); ok {
					if _, exists := request["ids"]; exists {
						mockService.AssertExpectations(t)
					}
				}
			}
		})
	}
}

func TestStaffHandler_SearchStaffByName(t *testing.T) {
	mockService := new(MockStaffService)
	handler := NewStaffHandler(mockService)

	mockStaff := []models.Staff{
		{ID: 1, Name: "田中太郎", Email: "tanaka@example.com"},
		{ID: 2, Name: "田中花子", Email: "tanaka-hanako@example.com"},
	}

	tests := []struct {
		name           string
		searchName     string
		mockReturn     []models.Staff
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功：名前検索",
			searchName:     "田中",
			mockReturn:     mockStaff,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "エラー：検索名が空",
			searchName:     "",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "エラー：データベースエラー",
			searchName:     "田中",
			mockReturn:     nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/staff/search?name="+tt.searchName, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.searchName != "" {
				mockService.On("SearchStaffByName", tt.searchName).Return(tt.mockReturn, tt.mockError)
			}

			err := handler.SearchStaffByName(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.searchName != "" {
				mockService.AssertExpectations(t)
			}
		})
	}
}