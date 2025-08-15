package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouterWithAccessLogger(t *testing.T) (r *gin.Engine, logPath string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Chạy trong thư mục tạm để log/access.log không đụng file thật
	tmp := t.TempDir()
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	t.Cleanup(func() { _ = os.Chdir(oldCwd) })

	// Gán vào biến return đã khai báo (không dùng :=) -> không thể dính SA4006
	r = gin.New()
	r.Use(AccessLogger())

	// Route không thuộc /api/v1 (để test skip)
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// POST /api/v1/echo: trả nguyên body
	r.POST("/api/v1/echo", func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", b)
	})

	// GET /api/v1/users: giả lập list endpoint
	r.GET("/api/v1/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"items": []int{1, 2}})
	})

	logPath = filepath.Join(tmp, "log", "access.log")
	return
}

func readAll(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		// file có thể chưa tạo cho đến khi middleware được khởi tạo (đã khởi tạo rồi)
		t.Fatalf("read file: %v", err)
	}
	return string(b)
}

func TestAccessLogger_SkipNonAPIRoute(t *testing.T) {
	r, logPath := setupRouterWithAccessLogger(t)

	// Gọi route không có /api/v1 -> middleware sẽ bỏ qua ghi nội dung
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	// File đã được tạo khi init middleware, nhưng không nên có nội dung log
	content := readAll(t, logPath)
	if strings.TrimSpace(content) != "" {
		t.Fatalf("expected empty log content for non-api route, got:\n%s", content)
	}
}

func TestAccessLogger_LogPostJSON(t *testing.T) {
	r, logPath := setupRouterWithAccessLogger(t)

	body := `{"a":1}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/echo?lang=en", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "req-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	content := readAll(t, logPath)

	// Kiểm tra các trường quan trọng có trong log
	if !strings.Contains(content, ". ID: req-1") {
		t.Fatalf("expected ID to be logged, got:\n%s", content)
	}
	if !strings.Contains(content, ". Method: POST") {
		t.Fatalf("expected Method POST in log")
	}
	if !strings.Contains(content, ". Path: /api/v1/echo?lang=en") {
		t.Fatalf("expected Path with query string in log")
	}
	// Request body pretty JSON
	if !strings.Contains(content, `. Request body (Content type: application/json):`) ||
		!strings.Contains(content, `"a": 1`) {
		t.Fatalf("expected pretty request body in log, got:\n%s", content)
	}
	// Response body pretty JSON
	if !strings.Contains(content, `. Response body (Content type: application/json):`) ||
		!strings.Contains(content, `"a": 1`) {
		t.Fatalf("expected pretty response body in log, got:\n%s", content)
	}
}

func TestAccessLogger_ListRouteHidesResponseBody(t *testing.T) {
	r, logPath := setupRouterWithAccessLogger(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=10", nil)
	req.Header.Set("X-Request-ID", "req-2")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	content := readAll(t, logPath)

	// Lấy đoạn log tương ứng với req-2 để kiểm tra nội dung
	idx := strings.LastIndex(content, ". ID: req-2")
	if idx == -1 {
		t.Fatalf("log segment for req-2 not found:\n%s", content)
	}
	seg := content[idx:] // từ ID: req-2 đến hết

	// Với GET list, formattedResp phải rỗng => không nên thấy "items"
	if strings.Contains(seg, `"items"`) {
		t.Fatalf("expected empty response body for list route, got:\n%s", seg)
	}
	if !strings.Contains(seg, ". Method: GET") {
		t.Fatalf("expected Method GET in segment")
	}
	if !strings.Contains(seg, ". Path: /api/v1/users?limit=10") {
		t.Fatalf("expected correct Path in segment")
	}
}
