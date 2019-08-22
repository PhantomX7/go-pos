package transaction

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type TransactionRepository interface {
	Insert(transaction *models.Transaction, tx *gorm.DB) error
	Update(transaction *models.Transaction, tx *gorm.DB) error
	Delete(transaction *models.Transaction, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.Transaction, error)
	FindByID(transactionID uint64) (*models.Transaction, error)
	FindByInvoiceID(invoiceID uint64) ([]models.Transaction, error)
	Count(config request_util.PaginationConfig) (int, error)
}
