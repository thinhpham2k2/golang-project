package controllers

import (
	"go-demo-gin/middlewares"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// Middleware stub: gắn localizer vào context để utils.LoadVariablesInContext không panic
func withStubLocalizer() gin.HandlerFunc {
	b := i18n.NewBundle(language.English) // không cần nạp file
	return func(c *gin.Context) {
		c.Set("localizer", i18n.NewLocalizer(b, "en"))
		c.Next()
	}
}

func TestLogin_BindingError_Returns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Dùng controller với svc = nil (OK vì binding fail sẽ return sớm)
	h := &AuthController{svc: nil}

	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/api/v1/authen/login", h.Login)

	// Gửi JSON lỗi để ShouldBind fail
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/authen/login", strings.NewReader(`{`)) // JSON lỗi
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	// Tuỳ utils.HandleBindError trả về cấu trúc JSON gì,
	// ở đây chỉ check có chữ "error" hoặc thông điệp parse JSON
	assert.Contains(t, strings.ToLower(w.Body.String()), "error")
}

func TestLogin_ServiceError_Returns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Dùng controller với svc = nil (OK vì binding fail sẽ return sớm)
	h := &AuthController{svc: nil}

	r := gin.New()
	r.Use(withStubLocalizer())
	r.Use(middlewares.ErrorHandler())
	r.POST("/api/v1/authen/login", h.Login)

	// Gửi JSON lỗi để ShouldBind fail
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/authen/login", strings.NewReader(`{}`)) // JSON lỗi
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	// Tuỳ utils.HandleBindError trả về cấu trúc JSON gì,
	// ở đây chỉ check có chữ "error" hoặc thông điệp parse JSON
	assert.Contains(t, strings.ToLower(w.Body.String()), "error")
}
