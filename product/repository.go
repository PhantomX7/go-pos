package product

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request"
)

type ProductRepository interface {
	Insert(product *models.Product) error
	Update(product *models.Product) error
	FindAll(config request.PaginationConfig) ([]models.Product, error)
	FindByID(productID int64) (models.Product, error)
	Count(config request.PaginationConfig) (int, error)
}
