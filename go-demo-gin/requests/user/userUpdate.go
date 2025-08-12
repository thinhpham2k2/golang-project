package user

import (
	"go-demo-gin/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// ✅ Hàm validate custom
func (u *UserUpdate) Validate(c *gin.Context, v *utils.Validator) map[string]string {
	validate := validator.New()
	validate.RegisterValidation("password", v.PasswordValidator)
	validate.RegisterValidation("birthday", v.BirthdayValidator)
	validate.RegisterValidation("hashed", v.HashedValidator)
	validate.RegisterValidation("role", v.RoleValidator)

	err := validate.Struct(u)
	if err == nil {
		return nil
	}

	errorsMap := make(map[string]string)
	for _, fe := range err.(validator.ValidationErrors) {
		// Lấy localizer cho i18n
		localizer := utils.LoadVariablesInContext(c)

		field := fe.Field()
		tag := fe.Tag()

		switch field {
		case "Pass":
			switch tag {
			case "required":
				errorsMap["password"] = utils.LoadI18nMessage(localizer, utils.PASSWORD_REQUIRE, nil)
			case "password":
				errorsMap["password"] = utils.LoadI18nMessage(localizer, utils.INVALID_PASSWORD, nil)
			case "hashed":
				errorsMap["password"] = utils.LoadI18nMessage(localizer, utils.PASSWORD_ENCRYPTION_FAIL, nil)
			}
		case "Role":
			switch tag {
			case "required":
				errorsMap["role"] = utils.LoadI18nMessage(localizer, utils.ROLE_REQUIRE, nil)
			case "role":
				errorsMap["role"] = utils.LoadI18nMessage(localizer, utils.INVALID_ROLE, nil)
			}
		case "Date":
			errorsMap["birthday"] = utils.LoadI18nMessage(localizer, utils.INVALID_BIRTHDAY, nil)
		default:
			errorsMap[field] = utils.LoadI18nMessage(localizer, utils.INVALID_VALUE, nil)
		}
	}
	return errorsMap
}
