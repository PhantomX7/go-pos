package invoice

import (
	"github.com/PhantomX7/go-pos/entity"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type InvoiceUsecase interface {
	Create(request request.InvoiceCreateRequest) (*models.Invoice, error)
	Update(invoiceID uint64, request request.InvoiceUpdateRequest) (*models.Invoice, error)
	Delete(invoiceID uint64) error
	Index(paginationConfig request.InvoicePaginationConfig) ([]entity.InvoiceDetail, response_util.PaginationMeta, error)
	Show(invoiceID uint64) (*entity.InvoiceDetail, error)
	SyncInvoice(invoiceID uint64) error
}
