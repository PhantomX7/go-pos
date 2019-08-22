package stock_adjustment

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type StockAdjustmentRepository interface {
	Insert(stockAdjustment *models.StockAdjustment, tx *gorm.DB) error
	Delete(stockAdjustment *models.StockAdjustment, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.StockAdjustment, error)
	FindByID(stockAdjustmentID uint64) (*models.StockAdjustment, error)
	FindByProductID(productID uint64) (*models.StockAdjustment, error)
	Count(config request_util.PaginationConfig) (int, error)
}
