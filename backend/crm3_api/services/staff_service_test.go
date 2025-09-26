package services

import (
	"crm3_api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple unit tests for staff service logic
func TestStaffServiceLogic(t *testing.T) {
	// Test staff creation request conversion
	phone := "090-1234-5678"
	storeID := uint(1)

	request := models.StaffCreateRequest{
		Name:     "田中太郎",
		Email:    "tanaka@example.com",
		Phone:    &phone,
		Position: "店長",
		StoreID:  &storeID,
		HireDate: "2023-01-01",
		Status:   "active",
	}

	staff := request.ToStaff()

	assert.Equal(t, request.Name, staff.Name)
	assert.Equal(t, request.Email, staff.Email)
	assert.Equal(t, request.Phone, staff.Phone)
	assert.Equal(t, request.Position, staff.Position)
	assert.Equal(t, request.StoreID, staff.StoreID)
	assert.Equal(t, request.HireDate, staff.HireDate)
	assert.Equal(t, request.Status, staff.Status)
}

func TestStaffUpdateRequestToStaff(t *testing.T) {
	// Test staff update request conversion
	phone := "090-5678-1234"
	storeID := uint(2)

	updateRequest := models.StaffUpdateRequest{
		Name:     "佐藤花子",
		Email:    "sato@example.com",
		Phone:    &phone,
		Position: "副店長",
		StoreID:  &storeID,
		HireDate: "2023-02-01",
		Status:   "inactive",
	}

	staff := updateRequest.ToStaff()

	assert.Equal(t, updateRequest.Name, staff.Name)
	assert.Equal(t, updateRequest.Email, staff.Email)
	assert.Equal(t, updateRequest.Phone, staff.Phone)
	assert.Equal(t, updateRequest.Position, staff.Position)
	assert.Equal(t, updateRequest.StoreID, staff.StoreID)
	assert.Equal(t, updateRequest.HireDate, staff.HireDate)
	assert.Equal(t, updateRequest.Status, staff.Status)
}

func TestStaffSearchParamsValidation(t *testing.T) {
	// Test search params structure
	name := "田中"
	email := "tanaka@example.com"
	position := "店長"
	storeID := uint(1)
	status := "active"

	params := models.StaffSearchParams{
		Name:     &name,
		Email:    &email,
		Position: &position,
		StoreID:  &storeID,
		Status:   &status,
		Page:     1,
		Limit:    10,
	}

	assert.Equal(t, name, *params.Name)
	assert.Equal(t, email, *params.Email)
	assert.Equal(t, position, *params.Position)
	assert.Equal(t, storeID, *params.StoreID)
	assert.Equal(t, status, *params.Status)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 10, params.Limit)
}

func TestPaginatedStaffResponseStructure(t *testing.T) {
	// Test paginated response structure
	staff := []models.Staff{
		{ID: 1, Name: "田中太郎", Status: "active"},
		{ID: 2, Name: "佐藤花子", Status: "inactive"},
	}

	response := models.PaginatedStaffResponse{
		Staff:      staff,
		Total:      2,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	assert.Equal(t, 2, len(response.Staff))
	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 1, response.TotalPages)
}

func TestStaffStatusCountsStructure(t *testing.T) {
	// Test status counts structure
	counts := models.StaffStatusCounts{
		Active:   5,
		Inactive: 2,
		OnLeave:  1,
		Total:    8,
	}

	assert.Equal(t, int64(5), counts.Active)
	assert.Equal(t, int64(2), counts.Inactive)
	assert.Equal(t, int64(1), counts.OnLeave)
	assert.Equal(t, int64(8), counts.Total)
}