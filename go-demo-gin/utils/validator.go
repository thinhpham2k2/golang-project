package utils

import (
	"context"
	"go-demo-gin/models"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ctxKeyUpdateID struct{}

func WithUpdateID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyUpdateID{}, id)
}

func UpdateIDFrom(ctx context.Context) (uint, bool) {
	if v := ctx.Value(ctxKeyUpdateID{}); v != nil {
		if id, ok := v.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

type Validator struct {
	db *gorm.DB
	v  *validator.Validate
}

func NewValidator(db *gorm.DB) *Validator {
	v := validator.New()
	val := &Validator{db: db, v: v}

	// các rule tĩnh của bạn
	_ = v.RegisterValidation("role", val.roleValidator)
	_ = v.RegisterValidation("hashed", val.hashedValidator)
	_ = v.RegisterValidation("password", val.passwordValidator)
	_ = v.RegisterValidation("username", val.usernameValidator)
	_ = v.RegisterValidation("birthday", val.birthdayValidator)

	// ✅ rule trùng username có context (timeout/cancel, dùng chung TX)
	_ = v.RegisterValidationCtx("duplicateUsername", val.duplicateUsernameCtx)

	return val
}

func (val *Validator) ValidateStructCtx(ctx context.Context, s any) map[string]string {
	cctx, cancel := context.WithTimeout(ctx, 700*time.Millisecond)
	defer cancel()

	if err := val.v.StructCtx(cctx, s); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			// Lấy localizer cho i18n
			localizer := LocalizerFrom(ctx)
			errorsMap := make(map[string]string)
			for _, fe := range verrs {
				field := fe.StructField()
				tag := fe.Tag()

				switch field {
				case "Username":
					switch tag {
					case "required":
						errorsMap["username"] = LoadI18nMessage(localizer, USERNAME_REQUIRE, nil)
					case "username":
						errorsMap["username"] = LoadI18nMessage(localizer, INVALID_USERNAME, nil)
					case "duplicateUsername":
						errorsMap["username"] = LoadI18nMessage(localizer, DUPLICATE_USERNAME, nil)
					}
				case "Pass":
					switch tag {
					case "required":
						errorsMap["password"] = LoadI18nMessage(localizer, PASSWORD_REQUIRE, nil)
					case "password":
						errorsMap["password"] = LoadI18nMessage(localizer, INVALID_PASSWORD, nil)
					case "hashed":
						errorsMap["password"] = LoadI18nMessage(localizer, PASSWORD_ENCRYPTION_FAIL, nil)
					}
				case "Role":
					switch tag {
					case "required":
						errorsMap["role"] = LoadI18nMessage(localizer, ROLE_REQUIRE, nil)
					case "role":
						errorsMap["role"] = LoadI18nMessage(localizer, INVALID_ROLE, nil)
					}
				case "Date":
					errorsMap["birthday"] = LoadI18nMessage(localizer, INVALID_BIRTHDAY, nil)
				default:
					errorsMap[field] = LoadI18nMessage(localizer, INVALID_VALUE, nil)
				}
			}

			return errorsMap
		}
	}
	return nil
}

func (v *Validator) roleValidator(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	switch models.Role(role) {
	case models.RoleAdmin, models.RoleStaff, models.RoleCustomer:
		return true
	default:
		return false
	}
}

func (v *Validator) hashedValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return err == nil
}

func (v *Validator) passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Regex: chỉ cho phép chữ thường, số, dấu chấm, gạch dưới; 3–24 ký tự
	re := regexp.MustCompile(`^[a-z0-9_.]{8,36}$`)
	return re.MatchString(password)
}

func (v *Validator) usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Regex: chỉ cho phép chữ thường, số, dấu chấm, gạch dưới; 3–24 ký tự
	re := regexp.MustCompile(`^[a-z0-9_.]{3,24}$`)
	return re.MatchString(username)
}

func (val *Validator) duplicateUsernameCtx(ctx context.Context, fl validator.FieldLevel) bool {
	username := fl.Field().String()

	q := val.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username)
	if currID, ok := UpdateIDFrom(ctx); ok { // 👈 lấy ID đã gắn
		q = q.Where("id <> ?", currID)
	}

	var count int64
	return q.Count(&count).Error == nil && count == 0
}

func (v *Validator) birthdayValidator(fl validator.FieldLevel) bool {
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
