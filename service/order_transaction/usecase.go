package order_transaction

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/order_transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type OrderTransactionUsecase interface {
	Create(request request.OrderTransactionCreateRequest) (*models.OrderTransaction, error)
	Update(orderTransactionID uint64, request request.OrderTransactionUpdateRequest) (*models.OrderTransaction, error)
	Delete(orderTransactionID uint64) error
	Index(paginationConfig request.OrderTransactionPaginationConfig) ([]models.OrderTransaction, response_util.PaginationMeta, error)
	Show(orderTransactionID uint64) (*models.OrderTransaction, error)
}
