package customer

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
)

type CustomerRepository interface {
	Insert(customer *models.Customer) error
	Update(customer *models.Customer) error
	Delete(customer *models.Customer) error
	FindAll(config request_util.PaginationConfig) ([]models.Customer, error)
	FindByID(customerID uint64) (models.Customer, error)
	Count(config request_util.PaginationConfig) (int, error)
}
