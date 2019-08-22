package models

import "time"

type StockAdjustment struct {
	ID          uint64    `json:"id"`
	ProductID   uint64    `json:"-"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}
