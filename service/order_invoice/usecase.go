package order_invoice

import (
	"github.com/PhantomX7/go-pos/entity"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/order_invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type OrderInvoiceUsecase interface {
	Create(request request.OrderInvoiceCreateRequest) (*models.OrderInvoice, error)
	Update(orderInvoiceID uint64, request request.OrderInvoiceUpdateRequest) (*models.OrderInvoice, error)
	Delete(orderInvoiceID uint64) error
	Index(paginationConfig request.OrderInvoicePaginationConfig) ([]entity.OrderInvoiceDetail, response_util.PaginationMeta, error)
	Show(orderInvoiceID uint64) (*entity.OrderInvoiceDetail, error)
	SyncOrderInvoice(orderInvoiceID uint64) error
}
