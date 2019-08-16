package models

import "time"

type Transaction struct {
	ID              int64     `json:"id"`
	InvoiceId       int64     `json:"-"`
	ProductId       int64     `json:"-"`
	StockMutationId int64     `json:"-"`
	CapitalPrice    float64   `json:"capital_price"`
	SellPrice       float64   `json:"sell_price"`
	Amount          int       `json:"amount"`
	TotalSellPrice  float64   `json:"total_sell_price"`
	Profit          float64   `json:"profit"`
	Date            time.Time `json:"date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"-"`
}
