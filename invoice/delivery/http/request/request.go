package request

import (
	"github.com/PhantomX7/go-pos/utils/request"
	"time"
)

// request related struct

type InvoiceCreateRequest struct {
	CustomerId    int64     `form:"customer_id" binding:"required,exist=customers.id"`
	Date          time.Time `form:"date" time_format:"2006-01-02"`
	PaymentStatus bool      `form:"payment_status"`
	PaymentType   string    `form:"payment_type"`
	Description   *string   `form:"description"`
}

type InvoiceUpdateRequest struct {
	Date          time.Time `form:"date" time_format:"2006-01-02"`
	PaymentStatus bool      `form:"payment_status"`
	PaymentType   string    `form:"payment_type"`
	Description   *string   `form:"description"`
}

type InvoicePaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request.SearchStruct
}

func NewInvoicePaginationConfig(conditions map[string][]string) InvoicePaginationConfig {
	filterable := map[string]string{
		"id":             request.IdType,
		"customer_id":    request.IdType,
		"date":           request.DateType,
		"payment_status": request.BoolType,
		"payment_type":   request.StringType,
		"description":    request.StringType,
	}

	invoicePaginationConfig := InvoicePaginationConfig{
		limit:        request.BuildLimit(conditions),
		offset:       request.BuildOffset(conditions),
		order:        request.BuildOrder(conditions),
		searchClause: request.BuildSearchClause(conditions, filterable),
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

func (i InvoicePaginationConfig) SearchClause() (res request.SearchStruct) {
	return i.searchClause
}
