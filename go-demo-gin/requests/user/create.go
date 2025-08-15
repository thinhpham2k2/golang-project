package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserCreate struct {
	Username string `json:"username" validate:"required,username,duplicateUsername"`
	Pass     string `json:"password" validate:"required,password,hashed" default:"12345678"`
	Name     string `json:"full_name"`
	Role     string `json:"role" validate:"required,role" default:"customer"`
	Date     string `json:"birthday" validate:"birthday" default:"2006-01-02"`
}

func (u *UserCreate) Password() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	return string(hash)
}

func (u *UserCreate) Birthday() *time.Time {
	birthday, _ := time.Parse("2006-01-02", u.Date)
	return &birthday
}
