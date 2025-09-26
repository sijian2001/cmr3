package services

import (
	"crm3_api/models"

	"gorm.io/gorm"
)

type StoreService struct {
	db *gorm.DB
}

func NewStoreService(db *gorm.DB) *StoreService {
	return &StoreService{db: db}
}

// CreateStore 店舗を作成
func (s *StoreService) CreateStore(store *models.Store) error {
	return s.db.Create(store).Error
}

// GetStoreByID IDで店舗を取得
func (s *StoreService) GetStoreByID(id uint) (*models.Store, error) {
	var store models.Store
	err := s.db.Preload("Manager").First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

// GetStoreByName 名前で店舗を取得
func (s *StoreService) GetStoreByName(name string) (*models.Store, error) {
	var store models.Store
	err := s.db.Preload("Manager").Where("name = ?", name).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

// UpdateStore 店舗を更新
func (s *StoreService) UpdateStore(store *models.Store) error {
	return s.db.Save(store).Error
}

// DeleteStore 店舗を削除
func (s *StoreService) DeleteStore(id uint) error {
	result := s.db.Delete(&models.Store{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// ListStores 店舗一覧を取得（検索とページネーション機能付き）
func (s *StoreService) ListStores(params models.StoreSearchParams) ([]models.Store, int64, error) {
	var stores []models.Store
	var total int64

	query := s.db.Model(&models.Store{}).Preload("Manager")

	// 検索条件を適用
	if params.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.ManagerID != nil {
		query = query.Where("manager_id = ?", *params.ManagerID)
	}

	// 総数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションを適用して店舗を取得
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

// UpdateStoreStatus 店舗ステータスを更新
func (s *StoreService) UpdateStoreStatus(id uint, status string) (*models.Store, error) {
	var store models.Store

	// 現在の店舗を取得
	if err := s.db.First(&store, id).Error; err != nil {
		return nil, err
	}

	// ステータス遷移の妥当性をチェック
	if err := store.UpdateStatus(status); err != nil {
		return nil, err
	}

	// 更新実行
	if err := s.db.Save(&store).Error; err != nil {
		return nil, err
	}

	// 管理者情報を含めて再取得
	if err := s.db.Preload("Manager").First(&store, id).Error; err != nil {
		return nil, err
	}

	return &store, nil
}

// GetStoresByStatus ステータス別に店舗を取得
func (s *StoreService) GetStoresByStatus(status string) ([]models.Store, error) {
	var stores []models.Store
	err := s.db.Preload("Manager").Where("status = ?", status).Order("created_at DESC").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

// GetActiveStores 営業中の店舗一覧を取得
func (s *StoreService) GetActiveStores() ([]models.Store, error) {
	return s.GetStoresByStatus("active")
}

// GetInactiveStores 休業中の店舗一覧を取得
func (s *StoreService) GetInactiveStores() ([]models.Store, error) {
	return s.GetStoresByStatus("inactive")
}

// GetMaintenanceStores メンテナンス中の店舗一覧を取得
func (s *StoreService) GetMaintenanceStores() ([]models.Store, error) {
	return s.GetStoresByStatus("maintenance")
}

// GetStoreStatusCounts 店舗ステータス別の件数を取得
func (s *StoreService) GetStoreStatusCounts() (*models.StoreStatusCounts, error) {
	var counts models.StoreStatusCounts

	// 総店舗数
	if err := s.db.Model(&models.Store{}).Count(&counts.Total).Error; err != nil {
		return nil, err
	}

	// 営業中の店舗数
	if err := s.db.Model(&models.Store{}).Where("status = ?", "active").Count(&counts.Active).Error; err != nil {
		return nil, err
	}

	// 休業中の店舗数
	if err := s.db.Model(&models.Store{}).Where("status = ?", "inactive").Count(&counts.Inactive).Error; err != nil {
		return nil, err
	}

	// メンテナンス中の店舗数
	if err := s.db.Model(&models.Store{}).Where("status = ?", "maintenance").Count(&counts.Maintenance).Error; err != nil {
		return nil, err
	}

	return &counts, nil
}

// GetStoresByManager 管理者IDで店舗一覧を取得
func (s *StoreService) GetStoresByManager(managerID uint) ([]models.Store, error) {
	var stores []models.Store
	err := s.db.Preload("Manager").Where("manager_id = ?", managerID).Order("created_at DESC").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

// GetStoresWithoutManager 管理者が設定されていない店舗一覧を取得
func (s *StoreService) GetStoresWithoutManager() ([]models.Store, error) {
	var stores []models.Store
	err := s.db.Where("manager_id IS NULL").Order("created_at DESC").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

// AssignManager 店舗に管理者を割り当て
func (s *StoreService) AssignManager(storeID, managerID uint) error {
	// 管理者の存在チェック
	var manager models.User
	if err := s.db.First(&manager, managerID).Error; err != nil {
		return err
	}

	// 店舗の存在チェック
	var store models.Store
	if err := s.db.First(&store, storeID).Error; err != nil {
		return err
	}

	// 管理者を割り当て
	store.ManagerID = &managerID
	return s.db.Save(&store).Error
}

// UnassignManager 店舗から管理者の割り当てを解除
func (s *StoreService) UnassignManager(storeID uint) error {
	var store models.Store
	if err := s.db.First(&store, storeID).Error; err != nil {
		return err
	}

	store.ManagerID = nil
	return s.db.Save(&store).Error
}

// BulkUpdateStatus 複数店舗のステータスを一括更新
func (s *StoreService) BulkUpdateStatus(storeIDs []uint, status string) error {
	// ステータスの妥当性をチェック
	validStatuses := map[string]bool{
		"active":      true,
		"inactive":    true,
		"maintenance": true,
	}
	if !validStatuses[status] {
		return gorm.ErrInvalidValue
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, id := range storeIDs {
		var store models.Store
		if err := tx.First(&store, id).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 各店舗のステータス遷移をチェック
		if err := store.ValidateStatusTransition(status); err != nil {
			tx.Rollback()
			return err
		}

		store.Status = status
		if err := tx.Save(&store).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// IsStoreNameAvailable 店舗名が利用可能かチェック
func (s *StoreService) IsStoreNameAvailable(name string, excludeID *uint) (bool, error) {
	var count int64
	query := s.db.Model(&models.Store{}).Where("name = ?", name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}