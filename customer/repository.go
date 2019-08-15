package customer

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request"
)

type CustomerRepository interface {
	Insert(customer *models.Customer) error
	Update(customer *models.Customer) error
	Delete(customer *models.Customer) error
	FindAll(config request.PaginationConfig) ([]models.Customer, error)
	FindByID(customerID int64) (models.Customer, error)
	Count(config request.PaginationConfig) (int, error)
}
