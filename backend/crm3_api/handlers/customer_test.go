package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"crm3_api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomerHandlerTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler *CustomerHandler
	echo    *echo.Echo
}

func (suite *CustomerHandlerTestSuite) SetupSuite() {
	// インメモリSQLiteデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // テスト中はログを無効化
	})
	suite.Require().NoError(err)

	// テーブル作成
	err = db.AutoMigrate(&models.Customer{})
	suite.Require().NoError(err)

	suite.db = db
	suite.handler = NewCustomerHandler(db)
	suite.echo = echo.New()
}

func (suite *CustomerHandlerTestSuite) SetupTest() {
	// 各テスト前にテーブルをクリア
	suite.db.Exec("DELETE FROM customers")
	suite.db.Exec("DELETE FROM sqlite_sequence WHERE name='customers'")
}

func (suite *CustomerHandlerTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *CustomerHandlerTestSuite) TestCreateCustomer_Success() {
	requestBody := models.CustomerCreateRequest{
		Name:    "テスト顧客",
		Email:   "test@example.com",
		Phone:   stringPtr("090-1234-5678"),
		Address: stringPtr("東京都渋谷区1-1-1"),
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	err := suite.handler.CreateCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)

	var response models.CustomerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("テスト顧客", response.Name)
	suite.Equal("test@example.com", response.Email)
	suite.Equal("090-1234-5678", *response.Phone)
	suite.Equal("東京都渋谷区1-1-1", *response.Address)
	suite.NotZero(response.ID)
}

func (suite *CustomerHandlerTestSuite) TestCreateCustomer_ValidationError() {
	requestBody := models.CustomerCreateRequest{
		Name:  "", // 必須フィールドが空
		Email: "invalid-email", // 無効なメール形式
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	err := suite.handler.CreateCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Validation Error", response["error"])
}

func (suite *CustomerHandlerTestSuite) TestCreateCustomer_DuplicateEmail() {
	// 最初の顧客を作成
	customer := models.Customer{
		Name:  "既存顧客",
		Email: "existing@example.com",
	}
	suite.db.Create(&customer)

	// 同じメールアドレスで新しい顧客を作成しようとする
	requestBody := models.CustomerCreateRequest{
		Name:  "新規顧客",
		Email: "existing@example.com", // 重複メール
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	err := suite.handler.CreateCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusConflict, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Conflict", response["error"])
	suite.Equal("Email address already exists", response["message"])
}

func (suite *CustomerHandlerTestSuite) TestGetCustomer_Success() {
	// テストデータ作成
	customer := models.Customer{
		Name:    "テスト顧客",
		Email:   "test@example.com",
		Phone:   stringPtr("090-1234-5678"),
		Address: stringPtr("東京都渋谷区1-1-1"),
	}
	suite.db.Create(&customer)

	req := httptest.NewRequest(http.MethodGet, "/api/customers/"+strconv.Itoa(int(customer.ID)), nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(customer.ID)))

	err := suite.handler.GetCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var response models.CustomerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(customer.ID, response.ID)
	suite.Equal("テスト顧客", response.Name)
	suite.Equal("test@example.com", response.Email)
}

func (suite *CustomerHandlerTestSuite) TestGetCustomer_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/api/customers/999", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := suite.handler.GetCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Not Found", response["error"])
}

func (suite *CustomerHandlerTestSuite) TestGetCustomer_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/customers/invalid", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := suite.handler.GetCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CustomerHandlerTestSuite) TestGetCustomers_Success() {
	// テストデータ作成
	customers := []models.Customer{
		{Name: "顧客A", Email: "a@example.com"},
		{Name: "顧客B", Email: "b@example.com"},
		{Name: "顧客C", Email: "c@example.com"},
	}
	for _, customer := range customers {
		suite.db.Create(&customer)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/customers?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	err := suite.handler.GetCustomers(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var response models.PaginatedCustomerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response.Customers, 3)
	suite.Equal(int64(3), response.Total)
	suite.Equal(1, response.Page)
	suite.Equal(10, response.Limit)
	suite.Equal(1, response.TotalPages)
}

func (suite *CustomerHandlerTestSuite) TestGetCustomers_WithSearch() {
	// テストデータ作成
	customers := []models.Customer{
		{Name: "山田太郎", Email: "yamada@example.com"},
		{Name: "田中花子", Email: "tanaka@example.com"},
		{Name: "佐藤次郎", Email: "sato@example.com"},
	}
	for _, customer := range customers {
		suite.db.Create(&customer)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/customers?name=山田", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	err := suite.handler.GetCustomers(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var response models.PaginatedCustomerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response.Customers, 1)
	suite.Equal("山田太郎", response.Customers[0].Name)
}

func (suite *CustomerHandlerTestSuite) TestGetCustomers_Pagination() {
	// テストデータ作成（15件）
	for i := 1; i <= 15; i++ {
		customer := models.Customer{
			Name:  "顧客" + strconv.Itoa(i),
			Email: "customer" + strconv.Itoa(i) + "@example.com",
		}
		suite.db.Create(&customer)
	}

	// 2ページ目、1ページ5件で取得
	req := httptest.NewRequest(http.MethodGet, "/api/customers?page=2&limit=5", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	err := suite.handler.GetCustomers(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var response models.PaginatedCustomerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response.Customers, 5)
	suite.Equal(int64(15), response.Total)
	suite.Equal(2, response.Page)
	suite.Equal(5, response.Limit)
	suite.Equal(3, response.TotalPages)
}

func (suite *CustomerHandlerTestSuite) TestUpdateCustomer_Success() {
	// テストデータ作成
	customer := models.Customer{
		Name:  "元の名前",
		Email: "old@example.com",
	}
	suite.db.Create(&customer)

	requestBody := models.CustomerUpdateRequest{
		Name:    "更新された名前",
		Email:   "updated@example.com",
		Phone:   stringPtr("090-1234-5678"),
		Address: stringPtr("東京都渋谷区1-1-1"),
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/customers/"+strconv.Itoa(int(customer.ID)), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(customer.ID)))

	err := suite.handler.UpdateCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var response models.CustomerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("更新された名前", response.Name)
	suite.Equal("updated@example.com", response.Email)
	suite.Equal("090-1234-5678", *response.Phone)
}

func (suite *CustomerHandlerTestSuite) TestUpdateCustomer_NotFound() {
	requestBody := models.CustomerUpdateRequest{
		Name:  "テスト",
		Email: "test@example.com",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/customers/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := suite.handler.UpdateCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *CustomerHandlerTestSuite) TestDeleteCustomer_Success() {
	// テストデータ作成
	customer := models.Customer{
		Name:  "削除される顧客",
		Email: "delete@example.com",
	}
	suite.db.Create(&customer)

	req := httptest.NewRequest(http.MethodDelete, "/api/customers/"+strconv.Itoa(int(customer.ID)), nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(customer.ID)))

	err := suite.handler.DeleteCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	// データベースから削除されていることを確認（ソフトデリート）
	var deletedCustomer models.Customer
	err = suite.db.First(&deletedCustomer, customer.ID).Error
	suite.Error(err) // レコードが見つからないエラー
}

func (suite *CustomerHandlerTestSuite) TestDeleteCustomer_NotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/api/customers/999", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := suite.handler.DeleteCustomer(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func TestCustomerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerHandlerTestSuite))
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}