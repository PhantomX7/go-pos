package models

import "time"

type Customer struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Address   *string    `json:"address"`
	Phone     *string    `json:"phone"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
