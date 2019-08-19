package customer

import (
	"github.com/PhantomX7/go-pos/service/customer/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type CustomerUsecase interface {
	Create(request request.CustomerCreateRequest) (models.Customer, error)
	Update(customerID uint64, request request.CustomerUpdateRequest) (models.Customer, error)
	Delete(customerID uint64) error
	Index(paginationConfig request.CustomerPaginationConfig) ([]models.Customer, response_util.PaginationMeta, error)
	Show(customerID uint64) (models.Customer, error)
}
