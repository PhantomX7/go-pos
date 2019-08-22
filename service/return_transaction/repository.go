package return_transaction

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type ReturnTransactionRepository interface {
	Insert(returnTransaction *models.ReturnTransaction, tx *gorm.DB) error
	Update(returnTransaction *models.ReturnTransaction, tx *gorm.DB) error
	Delete(returnTransaction *models.ReturnTransaction, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.ReturnTransaction, error)
	FindByID(returnTransactionID uint64) (*models.ReturnTransaction, error)
	FindByTransactionID(transactionID uint64) (*models.ReturnTransaction, error)
	FindByInvoiceID(invoiceID uint64) ([]models.ReturnTransaction, error)
	Count(config request_util.PaginationConfig) (int, error)
}
