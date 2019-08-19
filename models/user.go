package models

import (
	"time"
)

type User struct {
	ID        uint64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	RoleId    uint64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
