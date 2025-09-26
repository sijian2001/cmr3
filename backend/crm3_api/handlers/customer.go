package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"crm3_api/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CustomerHandler 顧客関連のハンドラー
type CustomerHandler struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewCustomerHandler 顧客ハンドラーのコンストラクタ
func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{
		db:        db,
		validator: validator.New(),
	}
}

// GetCustomers 顧客一覧を取得
func (h *CustomerHandler) GetCustomers(c echo.Context) error {
	var params models.CustomerSearchParams

	// クエリパラメータをバインド
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid query parameters",
			"details": err.Error(),
		})
	}

	// バリデーション
	if err := h.validator.Struct(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation Error",
			"message": "Invalid parameters",
			"details": err.Error(),
		})
	}

	// デフォルト値設定
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// クエリビルダー
	query := h.db.Model(&models.Customer{})

	// 検索条件追加
	if params.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	}
	if params.Email != "" {
		query = query.Where("email ILIKE ?", "%"+params.Email+"%")
	}
	if params.Phone != "" {
		query = query.Where("phone ILIKE ?", "%"+params.Phone+"%")
	}

	// 総件数取得
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to count customers",
		})
	}

	// ソート設定
	var orderBy string
	validSortColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
		"updated_at": true,
	}

	if validSortColumns[params.SortBy] {
		orderBy = params.SortBy
		if params.SortDesc {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}
	} else {
		orderBy = "id ASC"
	}

	// ページネーション
	offset := (params.Page - 1) * params.Limit
	var customers []models.Customer
	if err := query.Order(orderBy).Limit(params.Limit).Offset(offset).Find(&customers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to fetch customers",
		})
	}

	// レスポンス構築
	customerResponses := make([]models.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = customer.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	response := models.PaginatedCustomerResponse{
		Customers:  customerResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, response)
}

// GetCustomer 顧客詳細を取得
func (h *CustomerHandler) GetCustomer(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid customer ID",
		})
	}

	var customer models.Customer
	if err := h.db.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "Not Found",
				"message": "Customer not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to fetch customer",
		})
	}

	return c.JSON(http.StatusOK, customer.ToResponse())
}

// CreateCustomer 顧客を作成
func (h *CustomerHandler) CreateCustomer(c echo.Context) error {
	var req models.CustomerCreateRequest

	// リクエストボディをバインド
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// バリデーション
	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation Error",
			"message": "Invalid input data",
			"details": err.Error(),
		})
	}

	// メールアドレスの重複チェック
	var existingCustomer models.Customer
	if err := h.db.Where("email = ?", strings.ToLower(req.Email)).First(&existingCustomer).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"error":   "Conflict",
			"message": "Email address already exists",
		})
	}

	// 顧客モデルに変換
	customer := req.ToCustomer()

	// データベースに保存
	if err := h.db.Create(&customer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to create customer",
		})
	}

	return c.JSON(http.StatusCreated, customer.ToResponse())
}

// UpdateCustomer 顧客を更新
func (h *CustomerHandler) UpdateCustomer(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid customer ID",
		})
	}

	var req models.CustomerUpdateRequest

	// リクエストボディをバインド
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// バリデーション
	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation Error",
			"message": "Invalid input data",
			"details": err.Error(),
		})
	}

	// 既存の顧客を取得
	var customer models.Customer
	if err := h.db.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "Not Found",
				"message": "Customer not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to fetch customer",
		})
	}

	// メールアドレスの重複チェック（自分以外）
	if strings.ToLower(req.Email) != strings.ToLower(customer.Email) {
		var existingCustomer models.Customer
		if err := h.db.Where("email = ? AND id != ?", strings.ToLower(req.Email), uint(id)).First(&existingCustomer).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error":   "Conflict",
				"message": "Email address already exists",
			})
		}
	}

	// 更新データを適用
	req.ApplyToCustomer(&customer)

	// データベースを更新
	if err := h.db.Save(&customer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to update customer",
		})
	}

	return c.JSON(http.StatusOK, customer.ToResponse())
}

// DeleteCustomer 顧客を削除
func (h *CustomerHandler) DeleteCustomer(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid customer ID",
		})
	}

	// 既存の顧客を取得
	var customer models.Customer
	if err := h.db.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "Not Found",
				"message": "Customer not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to fetch customer",
		})
	}

	// ソフトデリート実行
	if err := h.db.Delete(&customer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Database Error",
			"message": "Failed to delete customer",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Customer deleted successfully",
	})
}