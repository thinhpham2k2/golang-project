package utils

import (
	"go-demo-gin/models"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Validator struct{ db *gorm.DB }

func NewValidator(db *gorm.DB) *Validator {
	return &Validator{db: db}
}

func (v *Validator) RoleValidator(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	switch models.Role(role) {
	case models.RoleAdmin, models.RoleStaff, models.RoleCustomer:
		return true
	default:
		return false
	}
}

func (v *Validator) HashedValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return err == nil
}

func (v *Validator) PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Regex: chỉ cho phép chữ thường, số, dấu chấm, gạch dưới; 3–24 ký tự
	re := regexp.MustCompile(`^[a-z0-9_.]{8,36}$`)
	return re.MatchString(password)
}

func (v *Validator) UsernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Regex: chỉ cho phép chữ thường, số, dấu chấm, gạch dưới; 3–24 ký tự
	re := regexp.MustCompile(`^[a-z0-9_.]{3,24}$`)
	return re.MatchString(username)
}

func (v *Validator) DuplicateUsernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var count int64
	if err := v.db.Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		// thận trọng: khi lỗi DB, coi như không hợp lệ (hoặc tuỳ policy)
		return false
	}
	return count == 0
}

func (v *Validator) BirthdayValidator(fl validator.FieldLevel) bool {
	birthdayStr := fl.Field().String()
	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		return false
	}

	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}

	return age >= 5 && age <= 100 // hoặc < 100 nếu bạn không cho tròn 100
}
