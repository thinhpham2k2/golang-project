package utils

import (
	"context"
	"go-demo-gin/pkg"
	errorResponse "go-demo-gin/responses/error"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

type txKey struct{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
func TxFrom(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
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
