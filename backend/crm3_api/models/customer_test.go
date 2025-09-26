package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerTableName(t *testing.T) {
	customer := Customer{}
	assert.Equal(t, "customers", customer.TableName())
}

func TestCustomerToResponse(t *testing.T) {
	customer := Customer{
		ID:      1,
		Name:    "テスト顧客",
		Email:   "test@example.com",
		Phone:   stringPtr("090-1234-5678"),
		Address: stringPtr("東京都渋谷区1-1-1"),
	}

	response := customer.ToResponse()

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "テスト顧客", response.Name)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "090-1234-5678", *response.Phone)
	assert.Equal(t, "東京都渋谷区1-1-1", *response.Address)
}

func TestCustomerCreateRequestToCustomer(t *testing.T) {
	req := CustomerCreateRequest{
		Name:    "新規顧客",
		Email:   "new@example.com",
		Phone:   stringPtr("090-5678-1234"),
		Address: stringPtr("大阪府大阪市1-1-1"),
	}

	customer := req.ToCustomer()

	assert.Equal(t, "新規顧客", customer.Name)
	assert.Equal(t, "new@example.com", customer.Email)
	assert.Equal(t, "090-5678-1234", *customer.Phone)
	assert.Equal(t, "大阪府大阪市1-1-1", *customer.Address)
	assert.Equal(t, uint(0), customer.ID) // IDは未設定
}

func TestCustomerCreateRequestToCustomerWithNilFields(t *testing.T) {
	req := CustomerCreateRequest{
		Name:  "最小顧客",
		Email: "minimal@example.com",
		// Phone と Address は nil
	}

	customer := req.ToCustomer()

	assert.Equal(t, "最小顧客", customer.Name)
	assert.Equal(t, "minimal@example.com", customer.Email)
	assert.Nil(t, customer.Phone)
	assert.Nil(t, customer.Address)
}

func TestCustomerUpdateRequestApplyToCustomer(t *testing.T) {
	// 既存の顧客
	customer := Customer{
		ID:      1,
		Name:    "元の名前",
		Email:   "old@example.com",
		Phone:   stringPtr("090-0000-0000"),
		Address: stringPtr("元の住所"),
	}

	// 更新リクエスト
	req := CustomerUpdateRequest{
		Name:    "更新された名前",
		Email:   "updated@example.com",
		Phone:   stringPtr("090-1111-1111"),
		Address: stringPtr("更新された住所"),
	}

	req.ApplyToCustomer(&customer)

	assert.Equal(t, uint(1), customer.ID) // IDは変更されない
	assert.Equal(t, "更新された名前", customer.Name)
	assert.Equal(t, "updated@example.com", customer.Email)
	assert.Equal(t, "090-1111-1111", *customer.Phone)
	assert.Equal(t, "更新された住所", *customer.Address)
}

func TestCustomerUpdateRequestApplyToCustomerWithNilFields(t *testing.T) {
	// 既存の顧客（全フィールドあり）
	customer := Customer{
		ID:      1,
		Name:    "元の名前",
		Email:   "old@example.com",
		Phone:   stringPtr("090-0000-0000"),
		Address: stringPtr("元の住所"),
	}

	// 更新リクエスト（一部フィールドをnilに）
	req := CustomerUpdateRequest{
		Name:  "更新された名前",
		Email: "updated@example.com",
		// Phone と Address は nil
	}

	req.ApplyToCustomer(&customer)

	assert.Equal(t, "更新された名前", customer.Name)
	assert.Equal(t, "updated@example.com", customer.Email)
	assert.Nil(t, customer.Phone)   // nilに更新される
	assert.Nil(t, customer.Address) // nilに更新される
}

func TestCustomerSearchParamsDefaults(t *testing.T) {
	params := CustomerSearchParams{}

	// デフォルト値が適切に設定されることをテスト
	assert.Equal(t, "", params.Name)
	assert.Equal(t, "", params.Email)
	assert.Equal(t, "", params.Phone)
	assert.Equal(t, 0, params.Page)  // structのゼロ値
	assert.Equal(t, 0, params.Limit) // structのゼロ値
	assert.Equal(t, "", params.SortBy)
	assert.False(t, params.SortDesc)
}

func TestPaginatedCustomerResponseStructure(t *testing.T) {
	customers := []CustomerResponse{
		{ID: 1, Name: "顧客1", Email: "customer1@example.com"},
		{ID: 2, Name: "顧客2", Email: "customer2@example.com"},
	}

	response := PaginatedCustomerResponse{
		Customers:  customers,
		Total:      10,
		Page:       2,
		Limit:      5,
		TotalPages: 2,
	}

	assert.Len(t, response.Customers, 2)
	assert.Equal(t, int64(10), response.Total)
	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 5, response.Limit)
	assert.Equal(t, 2, response.TotalPages)
	assert.Equal(t, "顧客1", response.Customers[0].Name)
	assert.Equal(t, "顧客2", response.Customers[1].Name)
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}