package routes

import (
	"go-demo-gin/initializers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mở DB sqlite in-memory cho test
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite memory: %v", err)
	}
	return db
}

// Tạo set "METHOD PATH" từ danh sách route
func routeSet(rs []gin.RouteInfo) map[string]bool {
	out := make(map[string]bool, len(rs))
	for _, r := range rs {
		out[r.Method+" "+r.Path] = true
	}
	return out
}

func TestSetupRoutes(t *testing.T) {
	t.Setenv("SECRET", "test-secret") // tránh lỗi thiếu SECRET khi khởi tạo AuthService
	gin.SetMode(gin.TestMode)

	r := gin.New()
	db := openTestDB(t)

	// KHỞI TẠO ROUTER (không được panic)
	SetupRoutes(r, db)

	got := routeSet(r.Routes())
	expected := []string{
		"GET /swagger/*any",

		"POST /api/v1/users",
		"GET /api/v1/users",
		"GET /api/v1/users/:id",
		"PUT /api/v1/users/:id",
		"DELETE /api/v1/users/:id",

		"POST /api/v1/authen/login",
	}

	for _, ep := range expected {
		assert.Truef(t, got[ep], "route not registered: %s", ep)
	}
}

func TestProtectedRoutes(t *testing.T) {
	t.Setenv("SECRET", "test-secret")
	gin.SetMode(gin.TestMode)

	r := gin.New()
	db := openTestDB(t)

	// Khởi tạo i18n như main()
	if err := initializers.LoadI18n(); err != nil {
		t.Fatalf("load i18n: %v", err)
	}

	SetupRoutes(r, db)

	// Thiếu Authorization -> middleware Authentication phải chặn (401)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
