package models

import "time"

type Product struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Pinyin          *string   `json:"pinyin"`
	ImageUrl        *string   `json:"image_url"`
	Stock           float64   `json:"stock"`
	Unit            string    `json:"unit"`
	UnitAmount      *int      `json:"unit_amount"`
	Description     *string   `json:"description"`
	CapitalPrice    float64   `json:"capital_price"`
	SellPriceCredit float64   `json:"sell_price_credit"`
	SellPriceCash   float64   `json:"sell_price_cash"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
