package transaction

import (
	"github.com/PhantomX7/go-pos/service/transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type TransactionUsecase interface {
	Create(request request.TransactionCreateRequest) (models.Transaction, error)
	Update(transactionID uint64, request request.TransactionUpdateRequest) (models.Transaction, error)
	Delete(transactionID uint64) error
	Index(paginationConfig request.TransactionPaginationConfig) ([]models.Transaction, response_util.PaginationMeta, error)
	Show(transactionID uint64) (models.Transaction, error)
}
