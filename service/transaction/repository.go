package transaction

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
)

type TransactionRepository interface {
	Insert(transaction *models.Transaction) error
	Update(transaction *models.Transaction) error
	Delete(transaction *models.Transaction) error
	FindAll(config request_util.PaginationConfig) ([]models.Transaction, error)
	FindByID(transactionID int64) (models.Transaction, error)
	Count(config request_util.PaginationConfig) (int, error)
}
