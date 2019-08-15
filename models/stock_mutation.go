package models

import "time"

const StockMutationIN = "IN"
const StockMutationOUT = "OUT"

type StockMutation struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
