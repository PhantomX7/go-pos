package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
	"time"
)

// request related struct

type OrderInvoiceCreateRequest struct {
	Date          time.Time `form:"date" time_format:"2006-01-02" binding:"required"`
	PaymentStatus bool      `form:"payment_status" binding:"required"`
	Description   *string   `form:"description"`
}

type OrderInvoiceUpdateRequest struct {
	Date          *time.Time `form:"date" time_format:"2006-01-02"`
	PaymentStatus *bool      `form:"payment_status"`
	Description   *string    `form:"description"`
}

type OrderInvoicePaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewOrderInvoicePaginationConfig(conditions map[string][]string) OrderInvoicePaginationConfig {
	filterable := map[string]string{
		"id":             request_util.IdType,
		"date":           request_util.DateType,
		"payment_status": request_util.BoolType,
	}

	orderInvoicePaginationConfig := OrderInvoicePaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
	}

	return orderInvoicePaginationConfig
}

func (i OrderInvoicePaginationConfig) Limit() (res int) {
	return i.limit
}

func (i OrderInvoicePaginationConfig) Order() string {
	return i.order
}

func (i OrderInvoicePaginationConfig) Offset() (res int) {
	return i.offset
}

func (i OrderInvoicePaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return i.searchClause
}
