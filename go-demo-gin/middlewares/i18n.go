package middlewares

import (
	"go-demo-gin/initializers"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		accept := c.GetHeader("Accept-Language")

		localizer := i18n.NewLocalizer(initializers.Bundle, lang, accept)

		// Gắn vào context
		c.Set("localizer", localizer)

		c.Next()
	}
}
