package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Name     sql.NullString
	Birthday *time.Time `gorm:"type:date"`
	Role     Role       `gorm:"type:varchar(20)"`
}
