package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
)

// request related struct

type TransactionCreateRequest struct {
	InvoiceId    uint64  `form:"invoice_id" binding:"required,exist=invoices.id"`
	ProductId    uint64  `form:"product_id" binding:"required,exist=products.id"`
	CapitalPrice float64 `form:"capital_price" binding:"required"`
	SellPrice    float64 `form:"sell_price" binding:"required"`
	Amount       float64 `form:"amount" binding:"required"`
}

type TransactionUpdateRequest struct {
	CapitalPrice float64 `form:"capital_price"`
	SellPrice    float64 `form:"sell_price"`
	Amount       float64 `form:"amount"`
}

type TransactionPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewTransactionPaginationConfig(conditions map[string][]string) TransactionPaginationConfig {
	filterable := map[string]string{
		"id":            request_util.IdType,
		"invoice_id":    request_util.IdType,
		"product_id":    request_util.IdType,
		"capital_price": request_util.NumberType,
		"sell_price":    request_util.NumberType,
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
