package models

import (
	"gorm.io/gorm"
	"time"
)

type Staff struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
	Email     string         `gorm:"size:320;uniqueIndex;not null" json:"email" validate:"required,email,max=320"`
	Phone     *string        `gorm:"size:20" json:"phone,omitempty" validate:"omitempty,max=20"`
	Position  string         `gorm:"size:50;not null" json:"position" validate:"required,max=50"`
	StoreID   *uint          `gorm:"index" json:"store_id,omitempty"`
	HireDate  string         `gorm:"type:date;not null" json:"hire_date" validate:"required"`
	Status    string         `gorm:"size:20;not null;default:active" json:"status" validate:"required,oneof=active inactive on_leave"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 関連
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

type StaffCreateRequest struct {
	Name     string  `json:"name" validate:"required,max=100"`
	Email    string  `json:"email" validate:"required,email,max=320"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Position string  `json:"position" validate:"required,max=50"`
	StoreID  *uint   `json:"store_id,omitempty"`
	HireDate string  `json:"hire_date" validate:"required"`
	Status   string  `json:"status" validate:"required,oneof=active inactive on_leave"`
}

type StaffUpdateRequest struct {
	Name     string  `json:"name" validate:"required,max=100"`
	Email    string  `json:"email" validate:"required,email,max=320"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Position string  `json:"position" validate:"required,max=50"`
	StoreID  *uint   `json:"store_id,omitempty"`
	HireDate string  `json:"hire_date" validate:"required"`
	Status   string  `json:"status" validate:"required,oneof=active inactive on_leave"`
}

type StaffSearchParams struct {
	Name     *string `query:"name"`
	Email    *string `query:"email"`
	Position *string `query:"position"`
	StoreID  *uint   `query:"store_id"`
	Status   *string `query:"status"`
	Page     int     `query:"page"`
	Limit    int     `query:"limit"`
}

type PaginatedStaffResponse struct {
	Staff      []Staff `json:"staff"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"totalPages"`
}

type StaffStatusCounts struct {
	Active    int64 `json:"active"`
	Inactive  int64 `json:"inactive"`
	OnLeave   int64 `json:"on_leave"`
	Total     int64 `json:"total"`
}

func (r *StaffCreateRequest) ToStaff() *Staff {
	return &Staff{
		Name:     r.Name,
		Email:    r.Email,
		Phone:    r.Phone,
		Position: r.Position,
		StoreID:  r.StoreID,
		HireDate: r.HireDate,
		Status:   r.Status,
	}
}

func (r *StaffUpdateRequest) ToStaff() *Staff {
	return &Staff{
		Name:     r.Name,
		Email:    r.Email,
		Phone:    r.Phone,
		Position: r.Position,
		StoreID:  r.StoreID,
		HireDate: r.HireDate,
		Status:   r.Status,
	}
}

// GetStatusLabel returns Japanese label for status
func (s *Staff) GetStatusLabel() string {
	statusLabels := map[string]string{
		"active":    "在籍",
		"inactive":  "休職",
		"on_leave":  "休暇中",
	}

	if label, exists := statusLabels[s.Status]; exists {
		return label
	}
	return s.Status
}

// GetStatusColor returns color code for UI display
func (s *Staff) GetStatusColor() string {
	statusColors := map[string]string{
		"active":   "success",
		"inactive": "warning",
		"on_leave": "info",
	}

	if color, exists := statusColors[s.Status]; exists {
		return color
	}
	return "default"
}

// Status checker methods
func (s *Staff) IsActive() bool {
	return s.Status == "active"
}

func (s *Staff) IsInactive() bool {
	return s.Status == "inactive"
}

func (s *Staff) IsOnLeave() bool {
	return s.Status == "on_leave"
}

// CanBeActivated checks if staff can be activated
func (s *Staff) CanBeActivated() bool {
	return s.Status == "inactive" || s.Status == "on_leave"
}

// CanBeDeactivated checks if staff can be deactivated
func (s *Staff) CanBeDeactivated() bool {
	return s.Status == "active"
}

// CanGoOnLeave checks if staff can go on leave
func (s *Staff) CanGoOnLeave() bool {
	return s.Status == "active"
}

// ValidateStatusTransition validates if status transition is allowed
func (s *Staff) ValidateStatusTransition(newStatus string) error {
	// Same status is not allowed
	if s.Status == newStatus {
		return gorm.ErrInvalidValue
	}

	// Valid statuses
	validStatuses := []string{"active", "inactive", "on_leave"}
	isValid := false
	for _, status := range validStatuses {
		if newStatus == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return gorm.ErrInvalidValue
	}

	// All transitions between valid statuses are allowed
	// Business rules can be added here if needed
	return nil
}

// UpdateStatus updates staff status with validation
func (s *Staff) UpdateStatus(newStatus string) error {
	if err := s.ValidateStatusTransition(newStatus); err != nil {
		return err
	}

	s.Status = newStatus
	return nil
}

// GetWorkingYears calculates years since hire date
func (s *Staff) GetWorkingYears() int {
	hireDate, err := time.Parse("2006-01-02", s.HireDate)
	if err != nil {
		return 0
	}

	return int(time.Since(hireDate).Hours() / 24 / 365)
}

// IsNewEmployee checks if staff is new (less than 6 months)
func (s *Staff) IsNewEmployee() bool {
	hireDate, err := time.Parse("2006-01-02", s.HireDate)
	if err != nil {
		return false
	}

	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	return hireDate.After(sixMonthsAgo)
}

// IsAssignedToStore checks if staff is assigned to any store
func (s *Staff) IsAssignedToStore() bool {
	return s.StoreID != nil && *s.StoreID > 0
}

// GetStoreName returns store name or default text
func (s *Staff) GetStoreName() string {
	if s.Store != nil {
		return s.Store.Name
	}
	return "未割り当て"
}

// BeforeCreate hook
func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	// Validate hire date is not in the future
	hireDate, err := time.Parse("2006-01-02", s.HireDate)
	if err != nil {
		return err
	}

	if hireDate.After(time.Now()) {
		return gorm.ErrInvalidValue
	}

	return nil
}

// BeforeUpdate hook
func (s *Staff) BeforeUpdate(tx *gorm.DB) error {
	// Validate hire date is not in the future
	hireDate, err := time.Parse("2006-01-02", s.HireDate)
	if err != nil {
		return err
	}

	if hireDate.After(time.Now()) {
		return gorm.ErrInvalidValue
	}

	return nil
}