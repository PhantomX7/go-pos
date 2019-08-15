package invoice

import (
	"github.com/PhantomX7/go-pos/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/invoice/delivery/http/response"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type InvoiceUsecase interface {
	Create(request request.InvoiceCreateRequest) (models.Invoice, error)
	Update(invoiceID int64, request request.InvoiceUpdateRequest) (models.Invoice, error)
	Delete(invoiceID int64) error
	Index(paginationConfig request.InvoicePaginationConfig) (response.InvoiceResponse, response_util.PaginationMeta, error)
	Show(invoiceID int64) (models.Invoice, error)
}
