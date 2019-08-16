package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
	"time"
)

// request related struct

type TransactionCreateRequest struct {
	InvoiceId    int64   `form:"invoice_id" binding:"required,exist=invoices.id"`
	ProductId    int64   `form:"product_id" binding:"required,exist=products.id"`
	CapitalPrice float64 `form:"capital_price" binding:"required"`
	SellPrice    float64 `form:"sell_price" binding:"required"`
	Amount       float64 `form:"amount" binding:"required"`
}

type TransactionUpdateRequest struct {
	Date          time.Time `form:"date" time_format:"2006-01-02"`
	PaymentStatus bool      `form:"payment_status"`
	PaymentType   string    `form:"payment_type"`
	Description   *string   `form:"description"`
}

type TransactionPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewTransactionPaginationConfig(conditions map[string][]string) TransactionPaginationConfig {
	filterable := map[string]string{
		"id":             request_util.IdType,
		"customer_id":    request_util.IdType,
		"date":           request_util.DateType,
		"payment_status": request_util.BoolType,
		"payment_type":   request_util.StringType,
		"description":    request_util.StringType,
	}

	transactionPaginationConfig := TransactionPaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
	}

	return transactionPaginationConfig
}

func (i TransactionPaginationConfig) Limit() (res int) {
	return i.limit
}

func (i TransactionPaginationConfig) Order() string {
	return i.order
}

func (i TransactionPaginationConfig) Offset() (res int) {
	return i.offset
}

func (i TransactionPaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return i.searchClause
}
