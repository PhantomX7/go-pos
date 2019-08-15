package models

import "time"

type Invoice struct {
	ID             int64      `json:"id"`
	ImageUrl       *string    `json:"image_url"`
	CustomerId     int64      `json:"-"`
	TotalCapital   float64    `json:"total_capital"`
	TotalSellPrice float64    `json:"total_sell_price"`
	TotalProfit    float64    `json:"total_profit"`
	Description    *string    `json:"description"`
	PaymentType    string     `json:"payment_type"`
	PaymentStatus  bool       `json:"payment_status"`
	Date           time.Time  `json:"date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `json:"-"`
}
