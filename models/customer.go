package models

import "time"

type Customer struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Address   *string    `json:"address"`
	Phone     *string    `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
