package models

import "time"

type Transaction struct {
	ID              uint64     `json:"id"`
	InvoiceId       uint64     `json:"-"`
	ProductId       uint64     `json:"-"`
	StockMutationId uint64     `json:"-"`
	CapitalPrice    float64   `json:"capital_price"`
	SellPrice       float64   `json:"sell_price"`
	Amount          float64       `json:"amount"`
	TotalSellPrice  float64   `json:"total_sell_price"`
	Profit          float64   `json:"profit"`
	Date            time.Time `json:"date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"-"`
}
