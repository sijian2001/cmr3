package handlers

import (
	"crm3_api/models"
	"crm3_api/services"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type StaffHandler struct {
	staffService services.StaffService
	validator    *validator.Validate
}

func NewStaffHandler(staffService services.StaffService) *StaffHandler {
	return &StaffHandler{
		staffService: staffService,
		validator:    validator.New(),
	}
}

// GetAllStaff retrieves all staff with optional filters and pagination
// GET /api/staff
func (h *StaffHandler) GetAllStaff(c echo.Context) error {
	var params models.StaffSearchParams

	// Parse query parameters
	if name := c.QueryParam("name"); name != "" {
		params.Name = &name
	}

	if email := c.QueryParam("email"); email != "" {
		params.Email = &email
	}

	if position := c.QueryParam("position"); position != "" {
		params.Position = &position
	}

	if storeIDStr := c.QueryParam("store_id"); storeIDStr != "" {
		if storeID, err := strconv.ParseUint(storeIDStr, 10, 32); err == nil {
			storeIDUint := uint(storeID)
			params.StoreID = &storeIDUint
		}
	}

	if status := c.QueryParam("status"); status != "" {
		params.Status = &status
	}

	// Parse pagination
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	response, err := h.staffService.GetAllStaff(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetStaffByID retrieves a specific staff member by ID
// GET /api/staff/:id
func (h *StaffHandler) GetStaffByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	staff, err := h.staffService.GetStaffByID(uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// CreateStaff creates a new staff member
// POST /api/staff
func (h *StaffHandler) CreateStaff(c echo.Context) error {
	var request models.StaffCreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	staff, err := h.staffService.CreateStaff(request)
	if err != nil {
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "store not found" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, staff)
}

// UpdateStaff updates an existing staff member
// PUT /api/staff/:id
func (h *StaffHandler) UpdateStaff(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	var request models.StaffUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	staff, err := h.staffService.UpdateStaff(uint(id), request)
	if err != nil {
		if err.Error() == "staff not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "store not found" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// DeleteStaff deletes a staff member
// DELETE /api/staff/:id
func (h *StaffHandler) DeleteStaff(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	err = h.staffService.DeleteStaff(uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Staff deleted successfully",
	})
}

// UpdateStaffStatus updates staff status
// POST /api/staff/:id/status
func (h *StaffHandler) UpdateStaffStatus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	var request struct {
		Status string `json:"status" validate:"required,oneof=active inactive on_leave"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	staff, err := h.staffService.UpdateStaffStatus(uint(id), request.Status)
	if err != nil {
		if err.Error() == "staff not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// AssignToStore assigns staff to a store
// POST /api/staff/:id/store
func (h *StaffHandler) AssignToStore(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	var request struct {
		StoreID uint `json:"store_id" validate:"required,min=1"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	staff, err := h.staffService.AssignToStore(uint(id), request.StoreID)
	if err != nil {
		if err.Error() == "staff not found" || err.Error() == "store not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// UnassignFromStore unassigns staff from store
// DELETE /api/staff/:id/store
func (h *StaffHandler) UnassignFromStore(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	staff, err := h.staffService.UnassignFromStore(uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// GetStaffStatusCounts returns counts of staff by status
// GET /api/staff/status-counts
func (h *StaffHandler) GetStaffStatusCounts(c echo.Context) error {
	counts, err := h.staffService.GetStaffStatusCounts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, counts)
}

// GetStaffByStore retrieves all staff assigned to a specific store
// GET /api/stores/:store_id/staff
func (h *StaffHandler) GetStaffByStore(c echo.Context) error {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid store ID",
		})
	}

	staff, err := h.staffService.GetStaffByStore(uint(storeID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// GetUnassignedStaff retrieves all staff not assigned to any store
// GET /api/staff/unassigned
func (h *StaffHandler) GetUnassignedStaff(c echo.Context) error {
	staff, err := h.staffService.GetUnassignedStaff()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

// BulkUpdateStatus updates status for multiple staff members
// POST /api/staff/bulk-update-status
func (h *StaffHandler) BulkUpdateStatus(c echo.Context) error {
	var request struct {
		IDs    []uint `json:"ids" validate:"required,min=1"`
		Status string `json:"status" validate:"required,oneof=active inactive on_leave"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	err := h.staffService.BulkUpdateStatus(request.IDs, request.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Staff status updated successfully",
	})
}

// BulkAssignToStore assigns multiple staff members to a store
// POST /api/staff/bulk-assign-store
func (h *StaffHandler) BulkAssignToStore(c echo.Context) error {
	var request struct {
		IDs     []uint `json:"ids" validate:"required,min=1"`
		StoreID uint   `json:"store_id" validate:"required,min=1"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	err := h.staffService.BulkAssignToStore(request.IDs, request.StoreID)
	if err != nil {
		if err.Error() == "store not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Staff assigned to store successfully",
	})
}

// SearchStaffByName searches staff by name
// GET /api/staff/search
func (h *StaffHandler) SearchStaffByName(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Name parameter is required",
		})
	}

	staff, err := h.staffService.SearchStaffByName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}