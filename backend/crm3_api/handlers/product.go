package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"crm3_api/models"
	"crm3_api/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productService services.ProductServiceInterface
	validate       *validator.Validate
}

func NewProductHandler(productService services.ProductServiceInterface, validate *validator.Validate) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		validate:       validate,
	}
}

// CreateProduct 製品を作成
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req models.ProductCreateRequest
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

	product := req.ToProduct()
	if err := h.productService.CreateProduct(&product); err != nil {
		if err == gorm.ErrDuplicatedKey {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "このSKUは既に存在します",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品の作成に失敗しました",
		})
	}

	return c.JSON(http.StatusCreated, product)
}

// GetProduct 製品を取得
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "製品が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct 製品を更新
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	var req models.ProductUpdateRequest
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

	product := req.ToProduct()
	product.ID = uint(id)

	if err := h.productService.UpdateProduct(&product); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "製品が見つかりません",
			})
		}
		if err == gorm.ErrDuplicatedKey {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "このSKUは既に存在します",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct 製品を削除
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "製品が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品の削除に失敗しました",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ListProducts 製品一覧を取得（ページネーションと検索機能付き）
func (h *ProductHandler) ListProducts(c echo.Context) error {
	var params models.ProductSearchParams

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

	products, total, err := h.productService.ListProducts(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品一覧の取得に失敗しました",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	response := models.PaginatedProductResponse{
		Products:   products,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, response)
}

// AdjustStock 在庫数量を調整
func (h *ProductHandler) AdjustStock(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDの形式が正しくありません",
		})
	}

	var req models.StockAdjustmentRequest
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

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "製品が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品の取得に失敗しました",
		})
	}

	if err := product.AdjustStock(req.Quantity, req.Operation); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "在庫調整に失敗しました: " + err.Error(),
		})
	}

	if err := h.productService.UpdateProduct(product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "在庫の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "在庫を調整しました",
		"product":        product,
		"stock_status":   product.GetStockStatus(),
		"stock_value":    product.GetStockValue(),
		"adjustment_log": fmt.Sprintf("操作: %s, 数量: %d, 理由: %s", req.Operation, req.Quantity, req.Reason),
	})
}

// GetStockSummary 在庫サマリーを取得
func (h *ProductHandler) GetStockSummary(c echo.Context) error {
	summary, err := h.productService.GetStockSummary()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "在庫サマリーの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, summary)
}

// SearchProductsBySKU SKUで製品を検索
func (h *ProductHandler) SearchProductsBySKU(c echo.Context) error {
	sku := c.QueryParam("sku")
	if sku == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "SKUパラメータが必要です",
		})
	}

	product, err := h.productService.GetProductBySKU(sku)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "該当するSKUの製品が見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "製品の検索に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, product)
}