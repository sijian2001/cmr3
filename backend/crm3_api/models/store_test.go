package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestStoreCreateRequest_ToStore(t *testing.T) {
	address := "東京都渋谷区"
	phone := "03-1234-5678"
	email := "test@example.com"
	managerID := uint(1)

	req := StoreCreateRequest{
		Name:      "テスト店舗",
		Address:   &address,
		Phone:     &phone,
		Email:     &email,
		Status:    "active",
		ManagerID: &managerID,
	}

	store := req.ToStore()

	assert.Equal(t, req.Name, store.Name)
	assert.Equal(t, req.Address, store.Address)
	assert.Equal(t, req.Phone, store.Phone)
	assert.Equal(t, req.Email, store.Email)
	assert.Equal(t, req.Status, store.Status)
	assert.Equal(t, req.ManagerID, store.ManagerID)
}

func TestStoreUpdateRequest_ToStore(t *testing.T) {
	address := "大阪府大阪市"
	phone := "06-1234-5678"
	email := "updated@example.com"
	managerID := uint(2)

	req := StoreUpdateRequest{
		Name:      "更新された店舗",
		Address:   &address,
		Phone:     &phone,
		Email:     &email,
		Status:    "inactive",
		ManagerID: &managerID,
	}

	store := req.ToStore()

	assert.Equal(t, req.Name, store.Name)
	assert.Equal(t, req.Address, store.Address)
	assert.Equal(t, req.Phone, store.Phone)
	assert.Equal(t, req.Email, store.Email)
	assert.Equal(t, req.Status, store.Status)
	assert.Equal(t, req.ManagerID, store.ManagerID)
}

func TestStore_GetStatusLabel(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		expectedLabel  string
	}{
		{
			name:          "営業中",
			status:        "active",
			expectedLabel: "営業中",
		},
		{
			name:          "休業中",
			status:        "inactive",
			expectedLabel: "休業中",
		},
		{
			name:          "メンテナンス中",
			status:        "maintenance",
			expectedLabel: "メンテナンス中",
		},
		{
			name:          "無効なステータス",
			status:        "unknown",
			expectedLabel: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.status}
			label := store.GetStatusLabel()
			assert.Equal(t, tt.expectedLabel, label)
		})
	}
}

func TestStore_StatusCheckers(t *testing.T) {
	tests := []struct {
		name               string
		status             string
		expectedActive     bool
		expectedInactive   bool
		expectedMaintenance bool
	}{
		{
			name:               "営業中",
			status:             "active",
			expectedActive:     true,
			expectedInactive:   false,
			expectedMaintenance: false,
		},
		{
			name:               "休業中",
			status:             "inactive",
			expectedActive:     false,
			expectedInactive:   true,
			expectedMaintenance: false,
		},
		{
			name:               "メンテナンス中",
			status:             "maintenance",
			expectedActive:     false,
			expectedInactive:   false,
			expectedMaintenance: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.status}

			assert.Equal(t, tt.expectedActive, store.IsActive())
			assert.Equal(t, tt.expectedInactive, store.IsInactive())
			assert.Equal(t, tt.expectedMaintenance, store.IsUnderMaintenance())
		})
	}
}

func TestStore_CanBeActivated(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "休業中から営業開始可能",
			status:   "inactive",
			expected: true,
		},
		{
			name:     "メンテナンス中から営業開始可能",
			status:   "maintenance",
			expected: true,
		},
		{
			name:     "営業中から営業開始不可",
			status:   "active",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.status}
			assert.Equal(t, tt.expected, store.CanBeActivated())
		})
	}
}

func TestStore_CanBeDeactivated(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "営業中から休業可能",
			status:   "active",
			expected: true,
		},
		{
			name:     "休業中から休業不可",
			status:   "inactive",
			expected: false,
		},
		{
			name:     "メンテナンス中から休業不可",
			status:   "maintenance",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.status}
			assert.Equal(t, tt.expected, store.CanBeDeactivated())
		})
	}
}

func TestStore_CanGoToMaintenance(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "営業中からメンテナンス可能",
			status:   "active",
			expected: true,
		},
		{
			name:     "休業中からメンテナンス可能",
			status:   "inactive",
			expected: true,
		},
		{
			name:     "メンテナンス中からメンテナンス不可",
			status:   "maintenance",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.status}
			assert.Equal(t, tt.expected, store.CanGoToMaintenance())
		})
	}
}

func TestStore_ValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name        string
		currentStatus string
		newStatus   string
		expectedErr error
	}{
		{
			name:          "休業中から営業中への遷移",
			currentStatus: "inactive",
			newStatus:     "active",
			expectedErr:   nil,
		},
		{
			name:          "メンテナンス中から営業中への遷移",
			currentStatus: "maintenance",
			newStatus:     "active",
			expectedErr:   nil,
		},
		{
			name:          "営業中から休業中への遷移",
			currentStatus: "active",
			newStatus:     "inactive",
			expectedErr:   nil,
		},
		{
			name:          "営業中からメンテナンス中への遷移",
			currentStatus: "active",
			newStatus:     "maintenance",
			expectedErr:   nil,
		},
		{
			name:          "休業中からメンテナンス中への遷移",
			currentStatus: "inactive",
			newStatus:     "maintenance",
			expectedErr:   nil,
		},
		{
			name:          "営業中から営業中への遷移（無効）",
			currentStatus: "active",
			newStatus:     "active",
			expectedErr:   gorm.ErrInvalidValue,
		},
		{
			name:          "休業中から休業中への遷移（無効）",
			currentStatus: "inactive",
			newStatus:     "inactive",
			expectedErr:   gorm.ErrInvalidValue,
		},
		{
			name:          "メンテナンス中からメンテナンス中への遷移（無効）",
			currentStatus: "maintenance",
			newStatus:     "maintenance",
			expectedErr:   gorm.ErrInvalidValue,
		},
		{
			name:          "無効なステータスへの遷移",
			currentStatus: "active",
			newStatus:     "unknown",
			expectedErr:   gorm.ErrInvalidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.currentStatus}
			err := store.ValidateStatusTransition(tt.newStatus)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestStore_UpdateStatus(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		newStatus     string
		expectedErr   error
		expectedStatus string
	}{
		{
			name:          "正常なステータス更新",
			currentStatus: "inactive",
			newStatus:     "active",
			expectedErr:   nil,
			expectedStatus: "active",
		},
		{
			name:          "無効なステータス更新",
			currentStatus: "active",
			newStatus:     "active",
			expectedErr:   gorm.ErrInvalidValue,
			expectedStatus: "active", // 変更されない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.currentStatus}
			err := store.UpdateStatus(tt.newStatus)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedStatus, store.Status)
		})
	}
}

func TestStore_GetStatusColor(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		expectedColor string
	}{
		{
			name:          "営業中",
			status:        "active",
			expectedColor: "success",
		},
		{
			name:          "休業中",
			status:        "inactive",
			expectedColor: "warning",
		},
		{
			name:          "メンテナンス中",
			status:        "maintenance",
			expectedColor: "error",
		},
		{
			name:          "無効なステータス",
			status:        "unknown",
			expectedColor: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{Status: tt.status}
			color := store.GetStatusColor()
			assert.Equal(t, tt.expectedColor, color)
		})
	}
}

func TestStoreStatusCounts_Structure(t *testing.T) {
	counts := StoreStatusCounts{
		Active:      5,
		Inactive:    2,
		Maintenance: 1,
		Total:       8,
	}

	assert.Equal(t, int64(5), counts.Active)
	assert.Equal(t, int64(2), counts.Inactive)
	assert.Equal(t, int64(1), counts.Maintenance)
	assert.Equal(t, int64(8), counts.Total)
}

func TestPaginatedStoreResponse_Structure(t *testing.T) {
	stores := []Store{
		{ID: 1, Name: "店舗1", Status: "active"},
		{ID: 2, Name: "店舗2", Status: "inactive"},
	}

	response := PaginatedStoreResponse{
		Stores:     stores,
		Total:      2,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	assert.Equal(t, 2, len(response.Stores))
	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 1, response.TotalPages)
}

// 統合テスト：複数の機能を組み合わせたテスト
func TestStore_CompleteWorkflow(t *testing.T) {
	// 店舗作成リクエストから店舗を作成
	address := "神奈川県横浜市"
	phone := "045-1234-5678"
	email := "test@example.com"
	managerID := uint(1)

	createReq := StoreCreateRequest{
		Name:      "テスト店舗",
		Address:   &address,
		Phone:     &phone,
		Email:     &email,
		Status:    "active",
		ManagerID: &managerID,
	}

	store := createReq.ToStore()

	// 初期状態の確認
	assert.Equal(t, "テスト店舗", store.Name)
	assert.Equal(t, "active", store.Status)
	assert.True(t, store.IsActive())
	assert.Equal(t, "営業中", store.GetStatusLabel())
	assert.Equal(t, "success", store.GetStatusColor())

	// ステータス遷移：営業中 -> 休業中
	err := store.UpdateStatus("inactive")
	assert.NoError(t, err)
	assert.Equal(t, "inactive", store.Status)
	assert.True(t, store.IsInactive())
	assert.Equal(t, "休業中", store.GetStatusLabel())
	assert.Equal(t, "warning", store.GetStatusColor())

	// ステータス遷移：休業中 -> メンテナンス中
	err = store.UpdateStatus("maintenance")
	assert.NoError(t, err)
	assert.Equal(t, "maintenance", store.Status)
	assert.True(t, store.IsUnderMaintenance())
	assert.Equal(t, "メンテナンス中", store.GetStatusLabel())
	assert.Equal(t, "error", store.GetStatusColor())

	// ステータス遷移：メンテナンス中 -> 営業中
	err = store.UpdateStatus("active")
	assert.NoError(t, err)
	assert.Equal(t, "active", store.Status)
	assert.True(t, store.IsActive())

	// 無効なステータス遷移のテスト
	err = store.UpdateStatus("active") // 営業中 -> 営業中は無効
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidValue, err)

	// 更新リクエストで店舗情報を更新
	newAddress := "愛知県名古屋市"
	updateReq := StoreUpdateRequest{
		Name:      "更新されたテスト店舗",
		Address:   &newAddress,
		Phone:     &phone,
		Email:     &email,
		Status:    "inactive",
		ManagerID: &managerID,
	}

	updatedStore := updateReq.ToStore()
	updatedStore.ID = store.ID

	assert.Equal(t, "更新されたテスト店舗", updatedStore.Name)
	assert.Equal(t, "愛知県名古屋市", *updatedStore.Address)
	assert.Equal(t, "inactive", updatedStore.Status)
	assert.True(t, updatedStore.IsInactive())
}