package models

import (
	"time"

	"gorm.io/gorm"
)

// Store 店舗モデル
type Store struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null" validate:"required,min=1,max=100"`
	Address   *string   `json:"address,omitempty" validate:"omitempty,max=500"`
	Phone     *string   `json:"phone,omitempty" validate:"omitempty,max=20"`
	Email     *string   `json:"email,omitempty" validate:"omitempty,email,max=320"`
	Status    string    `json:"status" gorm:"not null" validate:"required,oneof=active inactive maintenance"`
	ManagerID *uint     `json:"manager_id,omitempty" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 関連
	Manager *User `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
}

// StoreCreateRequest 店舗作成リクエスト
type StoreCreateRequest struct {
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Address   *string `json:"address,omitempty" validate:"omitempty,max=500"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email,max=320"`
	Status    string  `json:"status" validate:"required,oneof=active inactive maintenance"`
	ManagerID *uint   `json:"manager_id,omitempty"`
}

// StoreUpdateRequest 店舗更新リクエスト
type StoreUpdateRequest struct {
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Address   *string `json:"address,omitempty" validate:"omitempty,max=500"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email,max=320"`
	Status    string  `json:"status" validate:"required,oneof=active inactive maintenance"`
	ManagerID *uint   `json:"manager_id,omitempty"`
}

// StoreStatusUpdateRequest 店舗ステータス更新リクエスト
type StoreStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive maintenance"`
}

// StoreSearchParams 店舗検索パラメータ
type StoreSearchParams struct {
	Name     string `query:"name"`
	Status   string `query:"status"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	ManagerID *uint  `query:"manager_id"`
}

// PaginatedStoreResponse ページネーション付き店舗レスポンス
type PaginatedStoreResponse struct {
	Stores     []Store `json:"stores"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}

// StoreStatusCounts 店舗ステータス別件数
type StoreStatusCounts struct {
	Active      int64 `json:"active"`
	Inactive    int64 `json:"inactive"`
	Maintenance int64 `json:"maintenance"`
	Total       int64 `json:"total"`
}

// ToStore StoreCreateRequestをStoreに変換
func (req *StoreCreateRequest) ToStore() Store {
	return Store{
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		ManagerID: req.ManagerID,
	}
}

// ToStore StoreUpdateRequestをStoreに変換
func (req *StoreUpdateRequest) ToStore() Store {
	return Store{
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		ManagerID: req.ManagerID,
	}
}

// BeforeCreate 作成前のフック
func (s *Store) BeforeCreate(tx *gorm.DB) error {
	// 店舗名の重複チェック
	var count int64
	if err := tx.Model(&Store{}).Where("name = ?", s.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	// 管理者の存在チェック（指定されている場合）
	if s.ManagerID != nil {
		var manager User
		if err := tx.First(&manager, *s.ManagerID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return gorm.ErrInvalidValue
			}
			return err
		}
	}

	return nil
}

// BeforeUpdate 更新前のフック
func (s *Store) BeforeUpdate(tx *gorm.DB) error {
	// 店舗名の重複チェック（自分以外）
	var count int64
	if err := tx.Model(&Store{}).Where("name = ? AND id != ?", s.Name, s.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	// 管理者の存在チェック（指定されている場合）
	if s.ManagerID != nil {
		var manager User
		if err := tx.First(&manager, *s.ManagerID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return gorm.ErrInvalidValue
			}
			return err
		}
	}

	return nil
}

// GetStatusLabel ステータスラベルを取得
func (s *Store) GetStatusLabel() string {
	statusLabels := map[string]string{
		"active":      "営業中",
		"inactive":    "休業中",
		"maintenance": "メンテナンス中",
	}
	if label, exists := statusLabels[s.Status]; exists {
		return label
	}
	return s.Status
}

// IsActive 営業中かどうかを判定
func (s *Store) IsActive() bool {
	return s.Status == "active"
}

// IsInactive 休業中かどうかを判定
func (s *Store) IsInactive() bool {
	return s.Status == "inactive"
}

// IsUnderMaintenance メンテナンス中かどうかを判定
func (s *Store) IsUnderMaintenance() bool {
	return s.Status == "maintenance"
}

// CanBeActivated 営業開始可能かどうかを判定
func (s *Store) CanBeActivated() bool {
	return s.Status == "inactive" || s.Status == "maintenance"
}

// CanBeDeactivated 休業可能かどうかを判定
func (s *Store) CanBeDeactivated() bool {
	return s.Status == "active"
}

// CanGoToMaintenance メンテナンス状態に移行可能かどうかを判定
func (s *Store) CanGoToMaintenance() bool {
	return s.Status == "active" || s.Status == "inactive"
}

// ValidateStatusTransition ステータス遷移の妥当性をチェック
func (s *Store) ValidateStatusTransition(newStatus string) error {
	switch newStatus {
	case "active":
		if !s.CanBeActivated() {
			return gorm.ErrInvalidValue
		}
	case "inactive":
		if !s.CanBeDeactivated() {
			return gorm.ErrInvalidValue
		}
	case "maintenance":
		if !s.CanGoToMaintenance() {
			return gorm.ErrInvalidValue
		}
	default:
		return gorm.ErrInvalidValue
	}
	return nil
}

// UpdateStatus ステータスを更新
func (s *Store) UpdateStatus(newStatus string) error {
	if err := s.ValidateStatusTransition(newStatus); err != nil {
		return err
	}
	s.Status = newStatus
	return nil
}

// GetStatusColor ステータス色を取得
func (s *Store) GetStatusColor() string {
	statusColors := map[string]string{
		"active":      "success",
		"inactive":    "warning",
		"maintenance": "error",
	}
	if color, exists := statusColors[s.Status]; exists {
		return color
	}
	return "default"
}