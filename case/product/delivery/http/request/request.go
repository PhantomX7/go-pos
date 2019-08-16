package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
)

// request related struct

type ProductCreateRequest struct {
	Name            string  `form:"name" binding:"required,unique=products.name"`
	Pinyin          *string `form:"pinyin" `
	Stock           float64 `form:"stock" binding:"required,gte=0"`
	Unit            string  `form:"unit" binding:"required"`
	UnitAmount      *int    `form:"unit_amount" binding:"omitempty,gte=0"`
	Description     *string `form:"description"`
	CapitalPrice    float64 `form:"capital_price" binding:"required,gte=0"`
	SellPriceCredit float64 `form:"sell_price_credit" binding:"required,gte=0"`
	SellPriceCash   float64 `form:"sell_price_cash" binding:"required,gte=0"`
}

type ProductUpdateRequest struct {
	Name            *string  `form:"name" binding:"omitempty,unique=products.name"`
	Pinyin          *string  `form:"pinyin" `
	Stock           *float64 `form:"stock" binding:"omitempty,gte=0"`
	Unit            *string  `form:"unit"`
	UnitAmount      *int     `form:"unit_amount" binding:"omitempty,gte=0"`
	Description     *string  `form:"description"`
	CapitalPrice    *float64 `form:"capital_price" binding:"omitempty,gte=0"`
	SellPriceCredit *float64 `form:"sell_price_credit" binding:"omitempty,gte=0"`
	SellPriceCash   *float64 `form:"sell_price_cash" binding:"omitempty,gte=0"`
}

type ProductPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewProductPaginationConfig(conditions map[string][]string) ProductPaginationConfig {
	filterable := map[string]string{
		"id":                request_util.IdType,
		"name":              request_util.StringType,
		"pinyin":            request_util.StringType,
		"unit":              request_util.StringType,
		"unit_amount":       request_util.NumberType,
		"description":       request_util.StringType,
		"capital_price":     request_util.NumberType,
		"sell_price_credit": request_util.NumberType,
		"sell_price_cash":   request_util.NumberType,
		"created_at":        request_util.DateType,
	}

	productPaginationConfig := ProductPaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
	}

	return productPaginationConfig
}

func (p ProductPaginationConfig) Limit() (res int) {
	return p.limit
}

func (p ProductPaginationConfig) Order() string {
	return p.order
}

func (p ProductPaginationConfig) Offset() (res int) {
	return p.offset
}

func (p ProductPaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return p.searchClause
}
