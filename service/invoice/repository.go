package invoice

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
)

type InvoiceRepository interface {
	Insert(invoice *models.Invoice) error
	Update(invoice *models.Invoice) error
	Delete(invoice *models.Invoice) error
	FindAll(config request_util.PaginationConfig) ([]models.Invoice, error)
	FindByID(invoiceID int64) (models.Invoice, error)
	Count(config request_util.PaginationConfig) (int, error)
}
