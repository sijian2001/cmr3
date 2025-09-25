package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"crm3_api/dto"
	"crm3_api/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// MockDB はテスト用のモックデータベース
type MockDB struct {
	users []models.User
	nextID uint
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	if user, ok := value.(*models.User); ok {
		m.nextID++
		user.ID = m.nextID
		m.users = append(m.users, *user)
	}
	return &gorm.DB{}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	// モック実装 - 実際の検索は省略
	return &gorm.DB{}
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if _, ok := dest.(*models.User); ok {
		// エラーを返して「見つからない」を模擬
		return &gorm.DB{Error: gorm.ErrRecordNotFound}
	}
	return &gorm.DB{}
}

func (m *MockDB) Transaction(fc func(tx *gorm.DB) error) error {
	return fc(&gorm.DB{})
}

// TestAuthHandler_RegisterMock テスト用のハンドラーを直接テスト
func TestAuthHandler_RegisterMock(t *testing.T) {
	// テスト用にJWT_SECRETを設定
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	e := echo.New()

	// モックハンドラーを作成
	handler := &AuthHandler{
		DB:        nil, // ここではnilでよい（直接テストするため）
		Validator: validator.New(),
	}

	tests := []struct {
		name       string
		payload    dto.RegisterRequest
		wantStatus int
	}{
		{
			name: "Invalid email format",
			payload: dto.RegisterRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing name",
			payload: dto.RegisterRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			payload: dto.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonPayload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Register(c)
			if err != nil {
				t.Errorf("Register() error = %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("Register() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAuthHandler_LoginMock(t *testing.T) {
	// テスト用にJWT_SECRETを設定
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	e := echo.New()

	// モックハンドラーを作成
	handler := &AuthHandler{
		DB:        nil,
		Validator: validator.New(),
	}

	tests := []struct {
		name       string
		payload    dto.LoginRequest
		wantStatus int
	}{
		{
			name: "Invalid email format",
			payload: dto.LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing email",
			payload: dto.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			payload: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonPayload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Login(c)
			if err != nil {
				t.Errorf("Login() error = %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("Login() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAuthHandler_GetProfileMock(t *testing.T) {
	e := echo.New()

	handler := &AuthHandler{
		DB:        nil,
		Validator: validator.New(),
	}

	tests := []struct {
		name       string
		userID     interface{}
		wantStatus int
	}{
		{
			name:       "User not found (mocked)",
			userID:     uint(999),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", tt.userID)

			err := handler.GetProfile(c)
			if err != nil {
				t.Errorf("GetProfile() error = %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("GetProfile() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}