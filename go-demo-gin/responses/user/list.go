package user

import "time"

type UserList struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"full_name"`
	Role     string    `json:"role"`
	Birthday time.Time `json:"birthday"`
}
