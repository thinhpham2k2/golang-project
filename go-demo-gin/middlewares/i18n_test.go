package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-demo-gin/initializers"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// Khởi tạo bundle i18n tối thiểu cho test (không cần file JSON)
func setupTestBundle(t *testing.T) {
	t.Helper()
	b := i18n.NewBundle(language.English)
	// Thêm message trực tiếp (không đọc file)
	b.AddMessages(language.English, &i18n.Message{ID: "greet", Other: "Hello"})
	b.AddMessages(language.Vietnamese, &i18n.Message{ID: "greet", Other: "Xin chào"})
	initializers.Bundle = b
	t.Cleanup(func() { initializers.Bundle = nil })
}

// Handler dùng localizer đã được middleware gắn vào context
func greetHandler(c *gin.Context) {
	val, ok := c.Get("localizer")
	if !ok {
		c.String(http.StatusInternalServerError, "no localizer")
		return
	}
	loc, _ := val.(*i18n.Localizer)
	msg, err := loc.Localize(&i18n.LocalizeConfig{MessageID: "greet"})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, msg)
}

func TestI18n_LocalizerIsAttached(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestBundle(t)

	r := gin.New()
	r.Use(I18n())
	r.GET("/greet", greetHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/greet", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello", w.Body.String()) // mặc định English
}

func TestI18n_QueryParamOverridesAcceptLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestBundle(t)

	r := gin.New()
	r.Use(I18n())
	r.GET("/greet", greetHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/greet?lang=vi", nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Xin chào", w.Body.String())
}

func TestI18n_UsesAcceptLanguageWhenNoQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestBundle(t)

	r := gin.New()
	r.Use(I18n())
	r.GET("/greet", greetHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/greet", nil)
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en;q=0.8")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Xin chào", w.Body.String())
}

func TestI18n_DefaultsToEnglish(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestBundle(t)

	r := gin.New()
	r.Use(I18n())
	r.GET("/greet", greetHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/greet", nil) // không lang, không header
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello", w.Body.String())
}
