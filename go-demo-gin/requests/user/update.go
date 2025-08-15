package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUpdate struct {
	Pass string `json:"password" validate:"omitempty,password,hashed" default:"12345678"`
	Name string `json:"full_name"`
	Role string `json:"role" validate:"required,role" default:"customer"`
	Date string `json:"birthday" validate:"birthday" default:"2006-01-02"`
}

func (u *UserUpdate) Password() *string {
	if u.Pass == "" {
		return nil // Không hash, giữ nguyên password cũ
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	result := string(hash)
	return &result
}

func (u *UserUpdate) Birthday() *time.Time {
	birthday, _ := time.Parse("2006-01-02", u.Date)
	return &birthday
}
