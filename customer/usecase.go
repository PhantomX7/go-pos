package customer

import (
	"github.com/PhantomX7/go-pos/customer/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response"
)

type CustomerUsecase interface {
	Create(request request.CustomerCreateRequest) (models.Customer, error)
	Update(customerID int64, request request.CustomerUpdateRequest) (models.Customer, error)
	Delete(customerID int64) error
	Index(paginationConfig request.CustomerPaginationConfig) ([]models.Customer, response.PaginationMeta, error)
	Show(customerID int64) (models.Customer, error)
}
