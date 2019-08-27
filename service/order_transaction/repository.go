package order_transaction

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type OrderTransactionRepository interface {
	Insert(orderTransaction *models.OrderTransaction, tx *gorm.DB) error
	Update(orderTransaction *models.OrderTransaction, tx *gorm.DB) error
	Delete(orderTransaction *models.OrderTransaction, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.OrderTransaction, error)
	FindByID(orderTransactionID uint64) (*models.OrderTransaction, error)
	FindByOrderInvoiceID(orderInvoiceID uint64) ([]models.OrderTransaction, error)
	Count(config request_util.PaginationConfig) (int, error)
}
