package handlers

import (
	"math"
	"net/http"
	"strconv"

	"crm3_api/models"
	"crm3_api/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StoreHandler struct {
	storeService services.StoreServiceInterface
	validate     *validator.Validate
}

func NewStoreHandler(storeService services.StoreServiceInterface, validate *validator.Validate) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
		validate:     validate,
	}
}

// CreateStore 店舗を作成
func (h *StoreHandler) CreateStore(c echo.Context) error {
	var req models.StoreCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストの形式が正しくありません",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "バリデーションエラー: " + err.Error(),
		})
	}

	store := req.ToStore()
	if err := h.storeService.CreateStore(&store); err != nil {
		if err == gorm.ErrDuplicatedKey {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "この店舗名は既に存在します",
			})
		}
		if err == gorm.ErrInvalidValue {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "指定された管理者が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗の作成に失敗しました",
		})
	}

	return c.JSON(http.StatusCreated, store)
}

// GetStore 店舗を取得
func (h *StoreHandler) GetStore(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	store, err := h.storeService.GetStoreByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "店舗が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, store)
}

// UpdateStore 店舗を更新
func (h *StoreHandler) UpdateStore(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	var req models.StoreUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストの形式が正しくありません",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "バリデーションエラー: " + err.Error(),
		})
	}

	store := req.ToStore()
	store.ID = uint(id)

	if err := h.storeService.UpdateStore(&store); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "店舗が見つかりません",
			})
		}
		if err == gorm.ErrDuplicatedKey {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "この店舗名は既に存在します",
			})
		}
		if err == gorm.ErrInvalidValue {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "指定された管理者が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, store)
}

// DeleteStore 店舗を削除
func (h *StoreHandler) DeleteStore(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	if err := h.storeService.DeleteStore(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "店舗が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗の削除に失敗しました",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ListStores 店舗一覧を取得（ページネーションと検索機能付き）
func (h *StoreHandler) ListStores(c echo.Context) error {
	var params models.StoreSearchParams

	// クエリパラメータをバインド
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "クエリパラメータの形式が正しくありません",
		})
	}

	// デフォルト値設定
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	stores, total, err := h.storeService.ListStores(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗一覧の取得に失敗しました",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	response := models.PaginatedStoreResponse{
		Stores:     stores,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateStoreStatus 店舗ステータスを更新
func (h *StoreHandler) UpdateStoreStatus(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	var req models.StoreStatusUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストの形式が正しくありません",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "バリデーションエラー: " + err.Error(),
		})
	}

	store, err := h.storeService.UpdateStoreStatus(uint(id), req.Status)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "店舗が見つかりません",
			})
		}
		if err == gorm.ErrInvalidValue {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "無効なステータス遷移です",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ステータスの更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ステータスを更新しました",
		"store":   store,
	})
}

// GetStoresByStatus ステータス別に店舗一覧を取得
func (h *StoreHandler) GetStoresByStatus(c echo.Context) error {
	status := c.QueryParam("status")
	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ステータスパラメータが必要です",
		})
	}

	// ステータスの妥当性をチェック
	validStatuses := map[string]bool{
		"active":      true,
		"inactive":    true,
		"maintenance": true,
	}
	if !validStatuses[status] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なステータスです",
		})
	}

	stores, err := h.storeService.GetStoresByStatus(status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗一覧の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stores": stores,
		"count":  len(stores),
		"status": status,
	})
}

// GetStoreStatusCounts 店舗ステータス別件数を取得
func (h *StoreHandler) GetStoreStatusCounts(c echo.Context) error {
	counts, err := h.storeService.GetStoreStatusCounts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ステータス別件数の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, counts)
}

// AssignManager 店舗に管理者を割り当て
func (h *StoreHandler) AssignManager(c echo.Context) error {
	storeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "店舗IDの形式が正しくありません",
		})
	}

	var req struct {
		ManagerID uint `json:"manager_id" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストの形式が正しくありません",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "バリデーションエラー: " + err.Error(),
		})
	}

	if err := h.storeService.AssignManager(uint(storeID), req.ManagerID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "店舗または管理者が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "管理者の割り当てに失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "管理者を割り当てました",
	})
}

// UnassignManager 店舗から管理者の割り当てを解除
func (h *StoreHandler) UnassignManager(c echo.Context) error {
	storeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "店舗IDの形式が正しくありません",
		})
	}

	if err := h.storeService.UnassignManager(uint(storeID)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "店舗が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "管理者の割り当て解除に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "管理者の割り当てを解除しました",
	})
}

// GetStoresByManager 管理者IDで店舗一覧を取得
func (h *StoreHandler) GetStoresByManager(c echo.Context) error {
	managerID, err := strconv.ParseUint(c.QueryParam("manager_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "管理者IDの形式が正しくありません",
		})
	}

	stores, err := h.storeService.GetStoresByManager(uint(managerID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗一覧の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stores":     stores,
		"count":      len(stores),
		"manager_id": managerID,
	})
}

// BulkUpdateStatus 複数店舗のステータスを一括更新
func (h *StoreHandler) BulkUpdateStatus(c echo.Context) error {
	var req struct {
		StoreIDs []uint `json:"store_ids" validate:"required,min=1"`
		Status   string `json:"status" validate:"required,oneof=active inactive maintenance"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストの形式が正しくありません",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "バリデーションエラー: " + err.Error(),
		})
	}

	if err := h.storeService.BulkUpdateStatus(req.StoreIDs, req.Status); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "指定された店舗の一部が見つかりません",
			})
		}
		if err == gorm.ErrInvalidValue {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "無効なステータス遷移が含まれています",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ステータスの一括更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "ステータスを一括更新しました",
		"count":     len(req.StoreIDs),
		"status":    req.Status,
		"store_ids": req.StoreIDs,
	})
}

// CheckStoreName 店舗名の利用可能性をチェック
func (h *StoreHandler) CheckStoreName(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "店舗名パラメータが必要です",
		})
	}

	var excludeID *uint
	if idStr := c.QueryParam("exclude_id"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			uid := uint(id)
			excludeID = &uid
		}
	}

	available, err := h.storeService.IsStoreNameAvailable(name, excludeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "店舗名チェックに失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":      name,
		"available": available,
	})
}