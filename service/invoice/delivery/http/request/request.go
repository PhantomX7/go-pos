package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
	"time"
)

// request related struct

type InvoiceCreateRequest struct {
	CustomerId    uint64    `form:"customer_id" binding:"required,exist=customers.id"`
	Date          time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
	PaymentStatus bool      `form:"payment_status" binding:"required"`
	PaymentType   string    `form:"payment_type" binding:"required"`
	Description   *string   `form:"description"`
}

type InvoiceUpdateRequest struct {
	Date          *time.Time `form:"date" time_format:"2006-01-02"`
	PaymentStatus *bool      `form:"payment_status"`
	PaymentType   *string    `form:"payment_type"`
	Description   *string    `form:"description"`
}

type InvoicePaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewInvoicePaginationConfig(conditions map[string][]string) InvoicePaginationConfig {
	filterable := map[string]string{
		"id":             request_util.IdType,
		"customer_id":    request_util.IdType,
		"date":           request_util.DateType,
		"payment_status": request_util.BoolType,
		"payment_type":   request_util.StringType,
		"description":    request_util.StringType,
	}

	invoicePaginationConfig := InvoicePaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
	}

	return invoicePaginationConfig
}

func (i InvoicePaginationConfig) Limit() (res int) {
	return i.limit
}

func (i InvoicePaginationConfig) Order() string {
	return i.order
}

func (i InvoicePaginationConfig) Offset() (res int) {
	return i.offset
}

func (i InvoicePaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return i.searchClause
}
