package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
)

// request related struct

type CustomerCreateRequest struct {
	Name    string  `form:"name" binding:"required,unique=customers.name"`
	Address *string `form:"address" `
	Phone   *string `form:"phone"`
}

type CustomerUpdateRequest struct {
	Name    *string `form:"name" binding:"omitempty,unique=customers.name"`
	Address *string `form:"address" `
	Phone   *string `form:"phone"`
}

type CustomerPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewCustomerPaginationConfig(conditions map[string][]string) CustomerPaginationConfig {
	filterable := map[string]string{
		"id":      request_util.IdType,
		"name":    request_util.StringType,
		"address": request_util.StringType,
		"phone":   request_util.StringType,
	}

	customerPaginationConfig := CustomerPaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
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

func (c CustomerPaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return c.searchClause
}
