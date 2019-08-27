package order_invoice

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type OrderInvoiceRepository interface {
	Insert(orderInvoice *models.OrderInvoice, tx *gorm.DB) error
	Update(orderInvoice *models.OrderInvoice, tx *gorm.DB) error
	Delete(orderInvoice *models.OrderInvoice, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.OrderInvoice, error)
	FindByID(orderInvoiceID uint64) (*models.OrderInvoice, error)
	Count(config request_util.PaginationConfig) (int, error)
}
