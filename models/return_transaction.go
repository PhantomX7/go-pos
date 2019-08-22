package models

import "time"

type ReturnTransaction struct {
	ID              uint64    `json:"id"`
	TransactionID   uint64    `json:"-"`
	StockMutationID uint64    `json:"-"`
	InvoiceID       uint64    `json:"-"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"-"`
}
