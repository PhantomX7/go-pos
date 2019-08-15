package invoice

import (
	"github.com/PhantomX7/go-pos/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response"
)

type InvoiceUsecase interface {
	Create(request request.InvoiceCreateRequest) (models.Invoice, error)
	Update(invoiceID int64, request request.InvoiceUpdateRequest) (models.Invoice, error)
	Delete(invoiceID int64) error
	Index(paginationConfig request.InvoicePaginationConfig) ([]models.Invoice, response.PaginationMeta, error)
	Show(invoiceID int64) (models.Invoice, error)
}
