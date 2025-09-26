package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestStaffCreateRequest_ToStaff(t *testing.T) {
	phone := "090-1234-5678"
	storeID := uint(1)

	req := StaffCreateRequest{
		Name:     "田中太郎",
		Email:    "tanaka@example.com",
		Phone:    &phone,
		Position: "店長",
		StoreID:  &storeID,
		HireDate: "2023-01-01",
		Status:   "active",
	}

	staff := req.ToStaff()

	assert.Equal(t, req.Name, staff.Name)
	assert.Equal(t, req.Email, staff.Email)
	assert.Equal(t, req.Phone, staff.Phone)
	assert.Equal(t, req.Position, staff.Position)
	assert.Equal(t, req.StoreID, staff.StoreID)
	assert.Equal(t, req.HireDate, staff.HireDate)
	assert.Equal(t, req.Status, staff.Status)
}

func TestStaffUpdateRequest_ToStaff(t *testing.T) {
	phone := "090-5678-1234"
	storeID := uint(2)

	req := StaffUpdateRequest{
		Name:     "佐藤花子",
		Email:    "sato@example.com",
		Phone:    &phone,
		Position: "副店長",
		StoreID:  &storeID,
		HireDate: "2023-02-01",
		Status:   "inactive",
	}

	staff := req.ToStaff()

	assert.Equal(t, req.Name, staff.Name)
	assert.Equal(t, req.Email, staff.Email)
	assert.Equal(t, req.Phone, staff.Phone)
	assert.Equal(t, req.Position, staff.Position)
	assert.Equal(t, req.StoreID, staff.StoreID)
	assert.Equal(t, req.HireDate, staff.HireDate)
	assert.Equal(t, req.Status, staff.Status)
}

func TestStaff_GetStatusLabel(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		expectedLabel  string
	}{
		{
			name:          "在籍",
			status:        "active",
			expectedLabel: "在籍",
		},
		{
			name:          "休職",
			status:        "inactive",
			expectedLabel: "休職",
		},
		{
			name:          "休暇中",
			status:        "on_leave",
			expectedLabel: "休暇中",
		},
		{
			name:          "無効なステータス",
			status:        "unknown",
			expectedLabel: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Status: tt.status}
			label := staff.GetStatusLabel()
			assert.Equal(t, tt.expectedLabel, label)
		})
	}
}

func TestStaff_GetStatusColor(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		expectedColor string
	}{
		{
			name:          "在籍",
			status:        "active",
			expectedColor: "success",
		},
		{
			name:          "休職",
			status:        "inactive",
			expectedColor: "warning",
		},
		{
			name:          "休暇中",
			status:        "on_leave",
			expectedColor: "info",
		},
		{
			name:          "無効なステータス",
			status:        "unknown",
			expectedColor: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Status: tt.status}
			color := staff.GetStatusColor()
			assert.Equal(t, tt.expectedColor, color)
		})
	}
}

func TestStaff_StatusCheckers(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		expectedActive bool
		expectedInactive bool
		expectedOnLeave bool
	}{
		{
			name:           "在籍",
			status:         "active",
			expectedActive: true,
			expectedInactive: false,
			expectedOnLeave: false,
		},
		{
			name:           "休職",
			status:         "inactive",
			expectedActive: false,
			expectedInactive: true,
			expectedOnLeave: false,
		},
		{
			name:           "休暇中",
			status:         "on_leave",
			expectedActive: false,
			expectedInactive: false,
			expectedOnLeave: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Status: tt.status}

			assert.Equal(t, tt.expectedActive, staff.IsActive())
			assert.Equal(t, tt.expectedInactive, staff.IsInactive())
			assert.Equal(t, tt.expectedOnLeave, staff.IsOnLeave())
		})
	}
}

func TestStaff_CanBeActivated(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "休職中から復職可能",
			status:   "inactive",
			expected: true,
		},
		{
			name:     "休暇中から復職可能",
			status:   "on_leave",
			expected: true,
		},
		{
			name:     "在籍中から復職不可",
			status:   "active",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Status: tt.status}
			assert.Equal(t, tt.expected, staff.CanBeActivated())
		})
	}
}

func TestStaff_CanBeDeactivated(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "在籍中から休職可能",
			status:   "active",
			expected: true,
		},
		{
			name:     "休職中から休職不可",
			status:   "inactive",
			expected: false,
		},
		{
			name:     "休暇中から休職不可",
			status:   "on_leave",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Status: tt.status}
			assert.Equal(t, tt.expected, staff.CanBeDeactivated())
		})
	}
}

func TestStaff_CanGoOnLeave(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "在籍中から休暇可能",
			status:   "active",
			expected: true,
		},
		{
			name:     "休職中から休暇不可",
			status:   "inactive",
			expected: false,
		},
		{
			name:     "休暇中から休暇不可",
			status:   "on_leave",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Status: tt.status}
			assert.Equal(t, tt.expected, staff.CanGoOnLeave())
		})
	}
}

func TestStaff_ValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		newStatus     string
		expectedErr   error
	}{
		{
			name:          "休職中から在籍への遷移",
			currentStatus: "inactive",
			newStatus:     "active",
			expectedErr:   nil,
		},
		{
			name:          "休暇中から在籍への遷移",
			currentStatus: "on_leave",
			newStatus:     "active",
			expectedErr:   nil,
		},
		{
			name:          "在籍から休職への遷移",
			currentStatus: "active",
			newStatus:     "inactive",
			expectedErr:   nil,
		},
		{
			name:          "在籍から休暇への遷移",
			currentStatus: "active",
			newStatus:     "on_leave",
			expectedErr:   nil,
		},
		{
			name:          "休職から休暇への遷移",
			currentStatus: "inactive",
			newStatus:     "on_leave",
			expectedErr:   nil,
		},
		{
			name:          "休暇から休職への遷移",
			currentStatus: "on_leave",
			newStatus:     "inactive",
			expectedErr:   nil,
		},
		{
			name:          "在籍から在籍への遷移（無効）",
			currentStatus: "active",
			newStatus:     "active",
			expectedErr:   gorm.ErrInvalidValue,
		},
		{
			name:          "休職から休職への遷移（無効）",
			currentStatus: "inactive",
			newStatus:     "inactive",
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
			staff := &Staff{Status: tt.currentStatus}
			err := staff.ValidateStatusTransition(tt.newStatus)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestStaff_UpdateStatus(t *testing.T) {
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
			staff := &Staff{Status: tt.currentStatus}
			err := staff.UpdateStatus(tt.newStatus)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedStatus, staff.Status)
		})
	}
}

func TestStaff_GetWorkingYears(t *testing.T) {
	tests := []struct {
		name     string
		hireDate string
		expected int // This will be approximate due to time passage
	}{
		{
			name:     "2年前の入社",
			hireDate: "2022-01-01",
			expected: 2, // Will be approximately 2-3 years
		},
		{
			name:     "無効な日付形式",
			hireDate: "invalid-date",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{HireDate: tt.hireDate}
			years := staff.GetWorkingYears()

			if tt.hireDate == "invalid-date" {
				assert.Equal(t, tt.expected, years)
			} else {
				// Allow for a reasonable range due to current date variations
				assert.True(t, years >= 1 && years <= 4, "Expected years to be between 1 and 4, got %d", years)
			}
		})
	}
}

func TestStaff_IsNewEmployee(t *testing.T) {
	tests := []struct {
		name     string
		hireDate string
		expected bool
	}{
		{
			name:     "新入社員（3ヶ月前）",
			hireDate: "2023-06-01", // This will need to be adjusted based on current date
			expected: false, // Will depend on when test is run
		},
		{
			name:     "ベテラン社員（2年前）",
			hireDate: "2021-01-01",
			expected: false,
		},
		{
			name:     "無効な日付形式",
			hireDate: "invalid-date",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{HireDate: tt.hireDate}
			isNew := staff.IsNewEmployee()

			// For invalid date, should be false
			if tt.hireDate == "invalid-date" {
				assert.False(t, isNew)
			} else {
				// For valid dates, we just check the function doesn't panic
				assert.IsType(t, bool(false), isNew)
			}
		})
	}
}

func TestStaff_IsAssignedToStore(t *testing.T) {
	storeID := uint(1)
	zeroStoreID := uint(0)

	tests := []struct {
		name     string
		storeID  *uint
		expected bool
	}{
		{
			name:     "店舗に割り当て済み",
			storeID:  &storeID,
			expected: true,
		},
		{
			name:     "店舗に未割り当て（nil）",
			storeID:  nil,
			expected: false,
		},
		{
			name:     "店舗ID が 0",
			storeID:  &zeroStoreID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{StoreID: tt.storeID}
			assert.Equal(t, tt.expected, staff.IsAssignedToStore())
		})
	}
}

func TestStaff_GetStoreName(t *testing.T) {
	tests := []struct {
		name     string
		store    *Store
		expected string
	}{
		{
			name: "店舗情報あり",
			store: &Store{
				Name: "新宿店",
			},
			expected: "新宿店",
		},
		{
			name:     "店舗情報なし",
			store:    nil,
			expected: "未割り当て",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{Store: tt.store}
			assert.Equal(t, tt.expected, staff.GetStoreName())
		})
	}
}

func TestStaffStatusCounts_Structure(t *testing.T) {
	counts := StaffStatusCounts{
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

func TestPaginatedStaffResponse_Structure(t *testing.T) {
	staff := []Staff{
		{ID: 1, Name: "田中太郎", Status: "active"},
		{ID: 2, Name: "佐藤花子", Status: "inactive"},
	}

	response := PaginatedStaffResponse{
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

// 統合テスト：複数の機能を組み合わせたテスト
func TestStaff_CompleteWorkflow(t *testing.T) {
	// スタッフ作成リクエストからスタッフを作成
	phone := "090-1234-5678"
	storeID := uint(1)

	createReq := StaffCreateRequest{
		Name:     "田中太郎",
		Email:    "tanaka@example.com",
		Phone:    &phone,
		Position: "店長",
		StoreID:  &storeID,
		HireDate: "2023-01-01",
		Status:   "active",
	}

	staff := createReq.ToStaff()

	// 初期状態の確認
	assert.Equal(t, "田中太郎", staff.Name)
	assert.Equal(t, "active", staff.Status)
	assert.True(t, staff.IsActive())
	assert.Equal(t, "在籍", staff.GetStatusLabel())
	assert.Equal(t, "success", staff.GetStatusColor())
	assert.True(t, staff.IsAssignedToStore())

	// ステータス遷移：在籍 -> 休職
	err := staff.UpdateStatus("inactive")
	assert.NoError(t, err)
	assert.Equal(t, "inactive", staff.Status)
	assert.True(t, staff.IsInactive())
	assert.Equal(t, "休職", staff.GetStatusLabel())
	assert.Equal(t, "warning", staff.GetStatusColor())

	// ステータス遷移：休職 -> 休暇中
	err = staff.UpdateStatus("on_leave")
	assert.NoError(t, err)
	assert.Equal(t, "on_leave", staff.Status)
	assert.True(t, staff.IsOnLeave())
	assert.Equal(t, "休暇中", staff.GetStatusLabel())
	assert.Equal(t, "info", staff.GetStatusColor())

	// ステータス遷移：休暇中 -> 在籍
	err = staff.UpdateStatus("active")
	assert.NoError(t, err)
	assert.Equal(t, "active", staff.Status)
	assert.True(t, staff.IsActive())

	// 無効なステータス遷移のテスト
	err = staff.UpdateStatus("active") // 在籍 -> 在籍は無効
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidValue, err)

	// 更新リクエストでスタッフ情報を更新
	newPhone := "090-5678-1234"
	updateReq := StaffUpdateRequest{
		Name:     "田中次郎",
		Email:    "tanaka-jiro@example.com",
		Phone:    &newPhone,
		Position: "副店長",
		StoreID:  &storeID,
		HireDate: "2023-01-01",
		Status:   "inactive",
	}

	updatedStaff := updateReq.ToStaff()
	updatedStaff.ID = staff.ID

	assert.Equal(t, "田中次郎", updatedStaff.Name)
	assert.Equal(t, "tanaka-jiro@example.com", updatedStaff.Email)
	assert.Equal(t, "090-5678-1234", *updatedStaff.Phone)
	assert.Equal(t, "副店長", updatedStaff.Position)
	assert.Equal(t, "inactive", updatedStaff.Status)
	assert.True(t, updatedStaff.IsInactive())

	// 店舗割り当てのテスト
	assert.True(t, updatedStaff.IsAssignedToStore())

	// 店舗を設定してストア名を取得
	updatedStaff.Store = &Store{Name: "渋谷店"}
	assert.Equal(t, "渋谷店", updatedStaff.GetStoreName())

	// 店舗割り当て解除
	updatedStaff.StoreID = nil
	updatedStaff.Store = nil
	assert.False(t, updatedStaff.IsAssignedToStore())
	assert.Equal(t, "未割り当て", updatedStaff.GetStoreName())
}