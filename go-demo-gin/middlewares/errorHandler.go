package middlewares

import (
	errorResponse "go-demo-gin/responses/error"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors
		if len(errs) > 0 {
			err := errs[0].Err

			// Nếu là HTTPError, lấy status và message
			if httpErr, ok := err.(*errorResponse.HTTPError); ok {
				c.JSON(httpErr.StatusCode, httpErr.Message)
				return
			}

			// Lỗi thường
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}
	}
}
