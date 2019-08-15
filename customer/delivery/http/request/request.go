package request

import (
	"github.com/PhantomX7/go-pos/utils/request"
)

// request related struct

type CustomerCreateRequest struct {
	Name    string  `form:"name" binding:"required,unique=customers.name"`
	Address *string `form:"address" `
	Phone   *string `form:"phone"`
}

type CustomerUpdateRequest struct {
	Name    *string `form:"name" binding:"required,unique=customers.name"`
	Address *string `form:"address" `
	Phone   *string `form:"phone"`
}

type CustomerPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request.SearchStruct
}

func NewCustomerPaginationConfig(conditions map[string][]string) CustomerPaginationConfig {
	filterable := map[string]string{
		"id":      request.IdType,
		"name":    request.StringType,
		"address": request.StringType,
		"phone":   request.StringType,
	}

	customerPaginationConfig := CustomerPaginationConfig{
		limit:        request.BuildLimit(conditions),
		offset:       request.BuildOffset(conditions),
		order:        request.BuildOrder(conditions),
		searchClause: request.BuildSearchClause(conditions, filterable),
	}

	return customerPaginationConfig
}

func (c CustomerPaginationConfig) Limit() (res int) {
	return c.limit
}

func (c CustomerPaginationConfig) Order() string {
	return c.order
}

func (c CustomerPaginationConfig) Offset() (res int) {
	return c.offset
}

func (c CustomerPaginationConfig) SearchClause() (res request.SearchStruct) {
	return c.searchClause
}
