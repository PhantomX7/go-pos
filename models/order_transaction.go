package models

import "time"

type OrderTransaction struct {
	ID              uint64    `json:"id"`
	OrderInvoiceID  uint64    `json:"-"`
	ProductID       uint64    `json:"-"`
	StockMutationID uint64    `json:"-"`
	Amount          float64   `json:"amount"`
	BuyPrice        float64   `json:"buy_price"`
	TotalBuyPrice   float64   `json:"total_buy_price"`
	Date            time.Time `json:"date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"-"`
}
