package user

import (
	"time"
)

type UserDetail struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"full_name"`
	Role      string    `json:"role"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
