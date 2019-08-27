package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
)

// request related struct

type OrderTransactionCreateRequest struct {
	OrderInvoiceId uint64  `form:"order_invoice_id" binding:"required,exist=order_invoices.id"`
	ProductId      uint64  `form:"product_id" binding:"required,exist=products.id"`
	BuyPrice       float64 `form:"buy_price" binding:"required"`
	Amount         float64 `form:"amount" binding:"required"`
}

type OrderTransactionUpdateRequest struct {
	BuyPrice float64 `form:"buy_price" binding:"required,gte=0"`
	Amount   float64 `form:"amount" binding:"required,gte=0"`
}

type OrderTransactionPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewOrderTransactionPaginationConfig(conditions map[string][]string) OrderTransactionPaginationConfig {
	filterable := map[string]string{
		"id":               request_util.IdType,
		"order_invoice_id": request_util.IdType,
		"product_id":       request_util.IdType,
		"buy_price":        request_util.NumberType,
	}

	orderTransactionPaginationConfig := OrderTransactionPaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
	}

	return orderTransactionPaginationConfig
}

func (i OrderTransactionPaginationConfig) Limit() (res int) {
	return i.limit
}

func (i OrderTransactionPaginationConfig) Order() string {
	return i.order
}

func (i OrderTransactionPaginationConfig) Offset() (res int) {
	return i.offset
}

func (i OrderTransactionPaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return i.searchClause
}
