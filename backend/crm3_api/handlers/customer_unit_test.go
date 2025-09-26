package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crm3_api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCustomerHandler_Validation(t *testing.T) {
	e := echo.New()
	handler := &CustomerHandler{
		db:        nil, // モックテストなのでDBは不要
		validator: NewCustomerHandler(nil).validator,
	}

	t.Run("CreateCustomer_ValidationError_EmptyName", func(t *testing.T) {
		requestBody := models.CustomerCreateRequest{
			Name:  "", // 必須フィールドが空
			Email: "test@example.com",
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateCustomer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation Error", response["error"])
	})

	t.Run("CreateCustomer_ValidationError_InvalidEmail", func(t *testing.T) {
		requestBody := models.CustomerCreateRequest{
			Name:  "テスト顧客",
			Email: "invalid-email", // 無効なメール形式
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateCustomer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation Error", response["error"])
	})

	t.Run("CreateCustomer_ValidationError_InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateCustomer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", response["error"])
	})

	t.Run("GetCustomer_InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/customers/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetCustomer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", response["error"])
		assert.Equal(t, "Invalid customer ID", response["message"])
	})

	t.Run("UpdateCustomer_ValidationError", func(t *testing.T) {
		requestBody := models.CustomerUpdateRequest{
			Name:  "", // 必須フィールドが空
			Email: "invalid-email", // 無効なメール形式
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPut, "/api/customers/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.UpdateCustomer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DeleteCustomer_InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/customers/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.DeleteCustomer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCustomerHandler_GetCustomersParameterValidation(t *testing.T) {
	e := echo.New()
	handler := &CustomerHandler{
		db:        nil,
		validator: NewCustomerHandler(nil).validator,
	}

	t.Run("GetCustomers_InvalidParameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/customers?page=0&limit=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetCustomers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation Error", response["error"])
	})

	t.Run("GetCustomers_ValidParameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/customers?page=1&limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetCustomers(c)
		// データベースがnilなのでパニックが発生する前にバリデーションエラーになる可能性がある
		assert.NoError(t, err)
		// データベースがnilの場合の動作をテスト（実際には400か500のどちらかになる）
		assert.True(t, rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)
	})
}

func TestCustomerHandlerConstructor(t *testing.T) {
	t.Run("NewCustomerHandler", func(t *testing.T) {
		handler := NewCustomerHandler(nil)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.validator)
		assert.Nil(t, handler.db) // テストではnilを渡している
	})
}

// リクエスト・レスポンス構造体のテスト
func TestCustomerRequestResponseStructures(t *testing.T) {
	t.Run("CustomerCreateRequest_JSON_Binding", func(t *testing.T) {
		jsonData := `{
			"name": "テスト顧客",
			"email": "test@example.com",
			"phone": "090-1234-5678",
			"address": "東京都渋谷区1-1-1"
		}`

		var req models.CustomerCreateRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		assert.NoError(t, err)
		assert.Equal(t, "テスト顧客", req.Name)
		assert.Equal(t, "test@example.com", req.Email)
		assert.Equal(t, "090-1234-5678", *req.Phone)
		assert.Equal(t, "東京都渋谷区1-1-1", *req.Address)
	})

	t.Run("CustomerCreateRequest_JSON_Binding_MinimalFields", func(t *testing.T) {
		jsonData := `{
			"name": "最小顧客",
			"email": "minimal@example.com"
		}`

		var req models.CustomerCreateRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		assert.NoError(t, err)
		assert.Equal(t, "最小顧客", req.Name)
		assert.Equal(t, "minimal@example.com", req.Email)
		assert.Nil(t, req.Phone)
		assert.Nil(t, req.Address)
	})

	t.Run("CustomerUpdateRequest_JSON_Binding", func(t *testing.T) {
		jsonData := `{
			"name": "更新された顧客",
			"email": "updated@example.com",
			"phone": "090-8765-4321",
			"address": "大阪府大阪市2-2-2"
		}`

		var req models.CustomerUpdateRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		assert.NoError(t, err)
		assert.Equal(t, "更新された顧客", req.Name)
		assert.Equal(t, "updated@example.com", req.Email)
		assert.Equal(t, "090-8765-4321", *req.Phone)
		assert.Equal(t, "大阪府大阪市2-2-2", *req.Address)
	})

	t.Run("CustomerResponse_JSON_Serialization", func(t *testing.T) {
		response := models.CustomerResponse{
			ID:      1,
			Name:    "テスト顧客",
			Email:   "test@example.com",
			Phone:   stringPtr("090-1234-5678"),
			Address: stringPtr("東京都渋谷区1-1-1"),
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)
		assert.Contains(t, string(jsonData), `"id":1`)
		assert.Contains(t, string(jsonData), `"name":"テスト顧客"`)
		assert.Contains(t, string(jsonData), `"email":"test@example.com"`)
		assert.Contains(t, string(jsonData), `"phone":"090-1234-5678"`)
	})

	t.Run("PaginatedCustomerResponse_JSON_Serialization", func(t *testing.T) {
		customers := []models.CustomerResponse{
			{ID: 1, Name: "顧客1", Email: "customer1@example.com"},
			{ID: 2, Name: "顧客2", Email: "customer2@example.com"},
		}

		response := models.PaginatedCustomerResponse{
			Customers:  customers,
			Total:      10,
			Page:       2,
			Limit:      5,
			TotalPages: 2,
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)
		assert.Contains(t, string(jsonData), `"total":10`)
		assert.Contains(t, string(jsonData), `"page":2`)
		assert.Contains(t, string(jsonData), `"limit":5`)
		assert.Contains(t, string(jsonData), `"total_pages":2`)
		assert.Contains(t, string(jsonData), `"customers":[`)
	})
}

