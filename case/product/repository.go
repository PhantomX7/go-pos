package product

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
)

type ProductRepository interface {
	Insert(product *models.Product) error
	Update(product *models.Product) error
	FindAll(config request_util.PaginationConfig) ([]models.Product, error)
	FindByID(productID int64) (models.Product, error)
	Count(config request_util.PaginationConfig) (int, error)
}
