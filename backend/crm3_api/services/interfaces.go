package services

import "crm3_api/models"

// ProductServiceInterface defines the interface for product service
type ProductServiceInterface interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.Product, error)
	GetProductBySKU(sku string) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
	ListProducts(params models.ProductSearchParams) ([]models.Product, int64, error)
	GetStockSummary() (map[string]interface{}, error)
	GetProductsByStockStatus(status string) ([]models.Product, error)
	GetLowStockProducts() ([]models.Product, error)
	GetOutOfStockProducts() ([]models.Product, error)
	BulkUpdateStock(updates []struct {
		ID       uint `json:"id"`
		Quantity int  `json:"quantity"`
	}) error
}

// StoreServiceInterface defines the interface for store service
type StoreServiceInterface interface {
	CreateStore(store *models.Store) error
	GetStoreByID(id uint) (*models.Store, error)
	GetStoreByName(name string) (*models.Store, error)
	UpdateStore(store *models.Store) error
	DeleteStore(id uint) error
	ListStores(params models.StoreSearchParams) ([]models.Store, int64, error)
	UpdateStoreStatus(id uint, status string) (*models.Store, error)
	GetStoresByStatus(status string) ([]models.Store, error)
	GetActiveStores() ([]models.Store, error)
	GetInactiveStores() ([]models.Store, error)
	GetMaintenanceStores() ([]models.Store, error)
	GetStoreStatusCounts() (*models.StoreStatusCounts, error)
	GetStoresByManager(managerID uint) ([]models.Store, error)
	GetStoresWithoutManager() ([]models.Store, error)
	AssignManager(storeID, managerID uint) error
	UnassignManager(storeID uint) error
	BulkUpdateStatus(storeIDs []uint, status string) error
	IsStoreNameAvailable(name string, excludeID *uint) (bool, error)
}