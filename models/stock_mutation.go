package models

import "time"

type StockMutation struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
