package middlewares

import (
	"go-demo-gin/models"
	errorResponse "go-demo-gin/responses/error"
	"go-demo-gin/utils"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authentication(db *gorm.DB) func(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(allowedRoles ...models.Role) gin.HandlerFunc {
		return func(c *gin.Context) {
			// Lấy localizer cho i18n
			localizer := utils.LoadVariablesInContext(c)

			// 1. Lấy header Authorization
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.LoadI18nMessage(localizer, utils.INVALID_AUTHOR_HEADER, nil)})
				return
			}

			// 2. Cắt "Bearer " lấy token
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// 3. Parse và xác minh token
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(os.Getenv("SECRET")), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse.Error{
					Error: map[string]string{
						"message": err.Error(),
					},
				})
				return
			}

			// 4. Truy vấn thông tin user và phân quyền
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				username := claims["sub"]
				var user models.User
				result := db.First(&user, "username = ?", username)
				if result.Error != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse.Error{
						Error: map[string]string{
							"message": utils.LoadI18nMessage(localizer, utils.AUTHEN_REQUIRE, nil),
						},
					})
					return
				}

				c.Set("user", user) // Lưu thông tin user vào context

				if slices.Contains(allowedRoles, user.Role) {
					c.Next()
				} else {
					c.AbortWithStatusJSON(http.StatusForbidden, errorResponse.Error{
						Error: map[string]string{
							"message": utils.LoadI18nMessage(localizer, utils.PERMISSION_REQUIRE, nil),
						},
					})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.LoadI18nMessage(localizer, utils.INVALID_CLAIM, nil)})
				return
			}
		}
	}
}
