package product

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type ProductRepository interface {
	Insert(product *models.Product, tx *gorm.DB) error
	Update(product *models.Product, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.Product, error)
	FindByID(productID uint64) (models.Product, error)
	Count(config request_util.PaginationConfig) (int, error)
}
