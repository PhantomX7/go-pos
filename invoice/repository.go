package invoice

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request"
)

type InvoiceRepository interface {
	Insert(invoice *models.Invoice) error
	Update(invoice *models.Invoice) error
	Delete(invoice *models.Invoice) error
	FindAll(config request.PaginationConfig) ([]models.Invoice, error)
	FindByID(invoiceID int64) (models.Invoice, error)
	Count(config request.PaginationConfig) (int, error)
}
