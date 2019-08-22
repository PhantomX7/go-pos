package invoice

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type InvoiceRepository interface {
	Insert(invoice *models.Invoice, tx *gorm.DB) error
	Update(invoice *models.Invoice, tx *gorm.DB) error
	Delete(invoice *models.Invoice, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.Invoice, error)
	FindByID(invoiceID uint64) (*models.Invoice, error)
	Count(config request_util.PaginationConfig) (int, error)
}
