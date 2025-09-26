package services

import (
	"crm3_api/models"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type StaffService interface {
	GetAllStaff(params models.StaffSearchParams) (*models.PaginatedStaffResponse, error)
	GetStaffByID(id uint) (*models.Staff, error)
	CreateStaff(request models.StaffCreateRequest) (*models.Staff, error)
	UpdateStaff(id uint, request models.StaffUpdateRequest) (*models.Staff, error)
	DeleteStaff(id uint) error
	UpdateStaffStatus(id uint, status string) (*models.Staff, error)
	AssignToStore(id uint, storeID uint) (*models.Staff, error)
	UnassignFromStore(id uint) (*models.Staff, error)
	GetStaffStatusCounts() (*models.StaffStatusCounts, error)
	GetStaffByStore(storeID uint) ([]models.Staff, error)
	GetUnassignedStaff() ([]models.Staff, error)
	BulkUpdateStatus(ids []uint, status string) error
	BulkAssignToStore(ids []uint, storeID uint) error
	IsEmailAvailable(email string, excludeID *uint) (bool, error)
	SearchStaffByName(name string) ([]models.Staff, error)
}

type staffService struct {
	db *gorm.DB
}

func NewStaffService(db *gorm.DB) StaffService {
	return &staffService{db: db}
}

func (s *staffService) GetAllStaff(params models.StaffSearchParams) (*models.PaginatedStaffResponse, error) {
	var staff []models.Staff
	var total int64

	query := s.db.Model(&models.Staff{}).Preload("Store")

	// Apply filters
	if params.Name != nil && *params.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*params.Name+"%")
	}

	if params.Email != nil && *params.Email != "" {
		query = query.Where("email ILIKE ?", "%"+*params.Email+"%")
	}

	if params.Position != nil && *params.Position != "" {
		query = query.Where("position = ?", *params.Position)
	}

	if params.StoreID != nil {
		if *params.StoreID == 0 {
			// Search for unassigned staff
			query = query.Where("store_id IS NULL")
		} else {
			query = query.Where("store_id = ?", *params.StoreID)
		}
	}

	if params.Status != nil && *params.Status != "" {
		query = query.Where("status = ?", *params.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count staff: %w", err)
	}

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch staff: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &models.PaginatedStaffResponse{
		Staff:      staff,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *staffService) GetStaffByID(id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := s.db.Preload("Store").First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("staff not found")
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", err)
	}

	return &staff, nil
}

func (s *staffService) CreateStaff(request models.StaffCreateRequest) (*models.Staff, error) {
	// Check if email is already taken
	isAvailable, err := s.IsEmailAvailable(request.Email, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check email availability: %w", err)
	}
	if !isAvailable {
		return nil, fmt.Errorf("email already exists")
	}

	// Validate store exists if provided
	if request.StoreID != nil {
		var store models.Store
		if err := s.db.First(&store, *request.StoreID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("store not found")
			}
			return nil, fmt.Errorf("failed to validate store: %w", err)
		}
	}

	staff := request.ToStaff()

	if err := s.db.Create(staff).Error; err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	// Preload store information
	if err := s.db.Preload("Store").First(staff, staff.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch created staff: %w", err)
	}

	return staff, nil
}

func (s *staffService) UpdateStaff(id uint, request models.StaffUpdateRequest) (*models.Staff, error) {
	var existingStaff models.Staff
	if err := s.db.First(&existingStaff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("staff not found")
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", err)
	}

	// Check if email is available (excluding current staff)
	isAvailable, err := s.IsEmailAvailable(request.Email, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to check email availability: %w", err)
	}
	if !isAvailable {
		return nil, fmt.Errorf("email already exists")
	}

	// Validate store exists if provided
	if request.StoreID != nil {
		var store models.Store
		if err := s.db.First(&store, *request.StoreID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("store not found")
			}
			return nil, fmt.Errorf("failed to validate store: %w", err)
		}
	}

	// Update fields
	existingStaff.Name = request.Name
	existingStaff.Email = request.Email
	existingStaff.Phone = request.Phone
	existingStaff.Position = request.Position
	existingStaff.StoreID = request.StoreID
	existingStaff.HireDate = request.HireDate
	existingStaff.Status = request.Status

	if err := s.db.Save(&existingStaff).Error; err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	// Preload store information
	if err := s.db.Preload("Store").First(&existingStaff, existingStaff.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated staff: %w", err)
	}

	return &existingStaff, nil
}

func (s *staffService) DeleteStaff(id uint) error {
	var staff models.Staff
	if err := s.db.First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("staff not found")
		}
		return fmt.Errorf("failed to fetch staff: %w", err)
	}

	if err := s.db.Delete(&staff).Error; err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}

	return nil
}

func (s *staffService) UpdateStaffStatus(id uint, status string) (*models.Staff, error) {
	var staff models.Staff
	if err := s.db.First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("staff not found")
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", err)
	}

	// Validate status transition
	if err := staff.UpdateStatus(status); err != nil {
		return nil, fmt.Errorf("invalid status transition: %w", err)
	}

	if err := s.db.Save(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to update staff status: %w", err)
	}

	// Preload store information
	if err := s.db.Preload("Store").First(&staff, staff.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated staff: %w", err)
	}

	return &staff, nil
}

func (s *staffService) AssignToStore(id uint, storeID uint) (*models.Staff, error) {
	var staff models.Staff
	if err := s.db.First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("staff not found")
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", err)
	}

	// Validate store exists
	var store models.Store
	if err := s.db.First(&store, storeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("store not found")
		}
		return nil, fmt.Errorf("failed to validate store: %w", err)
	}

	staff.StoreID = &storeID

	if err := s.db.Save(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to assign staff to store: %w", err)
	}

	// Preload store information
	if err := s.db.Preload("Store").First(&staff, staff.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated staff: %w", err)
	}

	return &staff, nil
}

func (s *staffService) UnassignFromStore(id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := s.db.First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("staff not found")
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", err)
	}

	staff.StoreID = nil

	if err := s.db.Save(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to unassign staff from store: %w", err)
	}

	// Preload store information
	if err := s.db.Preload("Store").First(&staff, staff.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated staff: %w", err)
	}

	return &staff, nil
}

func (s *staffService) GetStaffStatusCounts() (*models.StaffStatusCounts, error) {
	var counts models.StaffStatusCounts

	// Count by status
	if err := s.db.Model(&models.Staff{}).Where("status = ?", "active").Count(&counts.Active).Error; err != nil {
		return nil, fmt.Errorf("failed to count active staff: %w", err)
	}

	if err := s.db.Model(&models.Staff{}).Where("status = ?", "inactive").Count(&counts.Inactive).Error; err != nil {
		return nil, fmt.Errorf("failed to count inactive staff: %w", err)
	}

	if err := s.db.Model(&models.Staff{}).Where("status = ?", "on_leave").Count(&counts.OnLeave).Error; err != nil {
		return nil, fmt.Errorf("failed to count on_leave staff: %w", err)
	}

	if err := s.db.Model(&models.Staff{}).Count(&counts.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total staff: %w", err)
	}

	return &counts, nil
}

func (s *staffService) GetStaffByStore(storeID uint) ([]models.Staff, error) {
	var staff []models.Staff
	if err := s.db.Where("store_id = ?", storeID).Preload("Store").Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch staff by store: %w", err)
	}

	return staff, nil
}

func (s *staffService) GetUnassignedStaff() ([]models.Staff, error) {
	var staff []models.Staff
	if err := s.db.Where("store_id IS NULL").Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch unassigned staff: %w", err)
	}

	return staff, nil
}

func (s *staffService) BulkUpdateStatus(ids []uint, status string) error {
	// Validate status
	validStatuses := []string{"active", "inactive", "on_leave"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid status: %s", status)
	}

	if err := s.db.Model(&models.Staff{}).Where("id IN ?", ids).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to bulk update staff status: %w", err)
	}

	return nil
}

func (s *staffService) BulkAssignToStore(ids []uint, storeID uint) error {
	// Validate store exists
	var store models.Store
	if err := s.db.First(&store, storeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("store not found")
		}
		return fmt.Errorf("failed to validate store: %w", err)
	}

	if err := s.db.Model(&models.Staff{}).Where("id IN ?", ids).Update("store_id", storeID).Error; err != nil {
		return fmt.Errorf("failed to bulk assign staff to store: %w", err)
	}

	return nil
}

func (s *staffService) IsEmailAvailable(email string, excludeID *uint) (bool, error) {
	query := s.db.Model(&models.Staff{}).Where("email = ?", email)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email availability: %w", err)
	}

	return count == 0, nil
}

func (s *staffService) SearchStaffByName(name string) ([]models.Staff, error) {
	var staff []models.Staff

	searchTerm := "%" + strings.ToLower(name) + "%"
	if err := s.db.Where("LOWER(name) LIKE ?", searchTerm).Preload("Store").Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to search staff by name: %w", err)
	}

	return staff, nil
}