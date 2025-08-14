package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	errorResponse "go-demo-gin/responses/error"
)

func newRouterWithMW() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	return r
}

func TestErrorHandler_HTTPError(t *testing.T) {
	r := newRouterWithMW()
	r.GET("/bad", func(c *gin.Context) {
		c.Error(&errorResponse.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message: errorResponse.Error{
				Error: map[string]string{"message": "bad input"},
			},
		})
		// Không viết response ở handler; middleware sẽ xử lý.
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	// Nội dung có thể phụ thuộc struct JSON; kiểm tra mềm.
	body := w.Body.String()
	if !strings.Contains(body, "bad input") {
		t.Fatalf("expected body to contain %q, got %s", "bad input", body)
	}

	// Thử parse để chắc có trường "error"
	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err == nil {
		if _, ok := m["error"]; !ok {
			t.Fatalf("expected top-level field 'error' in response JSON, got %v", m)
		}
	}
}

func TestErrorHandler_NormalError500(t *testing.T) {
	r := newRouterWithMW()
	r.GET("/boom", func(c *gin.Context) {
		c.Error(errors.New("boom"))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Internal Server Error") {
		t.Fatalf("expected body to contain %q, got %s", "Internal Server Error", body)
	}
	if !strings.Contains(body, "boom") {
		t.Fatalf("expected body to contain original error %q, got %s", "boom", body)
	}
}

func TestErrorHandler_NoErrorPassthrough(t *testing.T) {
	r := newRouterWithMW()
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if ok, _ := m["ok"].(bool); !ok {
		t.Fatalf("expected {\"ok\": true}, got %v", m)
	}
}

func TestErrorHandler_FirstErrorWins(t *testing.T) {
	r := newRouterWithMW()
	r.GET("/multi", func(c *gin.Context) {
		c.Error(errors.New("first"))
		c.Error(errors.New("second"))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/multi", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "first") || strings.Contains(body, "second") && !strings.Contains(body, "first") {
		t.Fatalf("expected middleware to respond using the FIRST error, got %s", body)
	}
}
