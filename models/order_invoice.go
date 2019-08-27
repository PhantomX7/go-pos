package models

import "time"

type OrderInvoice struct {
	ID            uint64    `json:"id"`
	TotalBuyPrice float64   `json:"total_buy_price"`
	PaymentStatus bool      `json:"payment_status"`
	Description   *string   `json:"description"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}
