package return_transaction

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/return_transaction/delivery/http/request"
)

type ReturnTransactionUsecase interface {
	Create(request request.ReturnTransactionCreateRequest) (*models.ReturnTransaction, error)
	Update(returnTransactionID uint64, request request.ReturnTransactionUpdateRequest) (*models.ReturnTransaction, error)
	Delete(returnTransactionID uint64) error
	Show(returnTransactionID uint64) (*models.ReturnTransaction, error)
}
