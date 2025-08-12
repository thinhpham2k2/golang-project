package utils

import (
	"go-demo-gin/models"
	"go-demo-gin/pkg"
	errorResponse "go-demo-gin/responses/error"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

func Log(c *gin.Context, level log.Level, message string) {
	// Lấy request id cho logging
	id, exists := c.Get("id")
	if !exists {
		id = "unknown"
	}

	// Lấy thông tin user cho logging
	val, exists := c.Get("user")  // exists: key có tồn tại không
	user, ok := val.(models.User) // ok: ép kiểu có thành công không
	if !exists || !ok {
		// Gán user mặc định (nếu models.User là struct)
		user = models.User{} // hoặc giá trị mặc định của bạn
	}

	entry := log.WithFields(log.Fields{
		"id":      id,
		"user_id": user.ID,
		"source":  "service",
	})

	switch level {
	case log.DebugLevel:
		entry.Debug(message)
	case log.InfoLevel:
		entry.Info(message)
	case log.WarnLevel:
		entry.Warn(message)
	case log.ErrorLevel:
		entry.Error(message)
	case log.FatalLevel:
		entry.Fatal(message)
	case log.TraceLevel:
		entry.Trace(message)
	default:
		entry.Print(message)
	}
}

func LoadVariablesInContext(c *gin.Context) *i18n.Localizer {
	// Lấy localizer cho i18n
	localizer := c.MustGet("localizer").(*i18n.Localizer)

	return localizer
}

func LoadI18nMessage(localizer *i18n.Localizer, message *i18n.Message, data map[string]any) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: message,
		TemplateData:   data,
		PluralCount:    -1,
	})
	if err != nil {
		return message.Other
	}
	return msg
}

func Paginate(pagination *pkg.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	totalPages := int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func HandleBindError(c *gin.Context, err error) {
	c.Error(&errorResponse.HTTPError{
		StatusCode: http.StatusBadRequest,
		Message: errorResponse.Error{
			Error: map[string]string{
				"message": err.Error(),
			},
		}})
}

func HandleValidationError(c *gin.Context, err map[string]string) {
	c.Error(&errorResponse.HTTPError{
		StatusCode: http.StatusBadRequest,
		Message: errorResponse.Error{
			Error: err,
		}})
}

func HandleServiceError(c *gin.Context, status int, msg string) {
	c.Error(&errorResponse.HTTPError{
		StatusCode: status,
		Message: errorResponse.Error{
			Error: map[string]string{
				"message": msg,
			},
		}})
}
