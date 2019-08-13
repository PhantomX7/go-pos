package request

import (
	"github.com/PhantomX7/go-pos/utils/request"
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
	searchClause request.SearchStruct
}

func NewProductPaginationConfig(conditions map[string][]string) ProductPaginationConfig {
	filterable := map[string]string{
		"id":                request.IdType,
		"name":              request.StringType,
		"pinyin":            request.StringType,
		"unit":              request.StringType,
		"unit_amount":       request.NumberType,
		"description":       request.StringType,
		"capital_price":     request.NumberType,
		"sell_price_credit": request.NumberType,
		"sell_price_cash":   request.NumberType,
		"created_at":        request.DateType,
	}

	userPaginationConfig := ProductPaginationConfig{
		limit:        request.BuildLimit(conditions),
		offset:       request.BuildOffset(conditions),
		order:        request.BuildOrder(conditions),
		searchClause: request.BuildSearchClause(conditions, filterable),
	}

	return userPaginationConfig
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

func (p ProductPaginationConfig) SearchClause() (res request.SearchStruct) {
	return p.searchClause
}
