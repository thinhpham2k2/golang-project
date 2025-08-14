package middlewares

import (
	"bytes"
	"encoding/json"
	"go-demo-gin/utils"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AccessLogger() gin.HandlerFunc {
	// Đường dẫn lưu trữ nhật kí
	logFilePath := "log/access.log"
	// Tạo thư mục nếu chưa có
	if err := os.MkdirAll(getDir(logFilePath), os.ModePerm); err != nil {
		log.Fatalf("Không thể tạo thư mục log: %v", err)
	}

	// Mở file log
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Không thể mở file log: %v", err)
	}
	logger := log.New(f, "", log.LstdFlags)

	return func(c *gin.Context) {
		if !strings.Contains(strings.ToLower(c.Request.RequestURI), "/api/v1") {
			c.Next()
			return
		}

		// Tạo ID cho lượt ghi nhật kí (logging)
		id := c.GetHeader("X-Request-ID")

		// Nếu không có hoặc giá trị không hợp lệ → server tự sinh
		if strings.TrimSpace(id) == "" {
			id = strconv.FormatInt(time.Now().UnixMilli(), 10) // time.Now().UnixNano()
		}

		start := time.Now()

		// Ghi lại request body
		var requestBody []byte
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = bodyBytes
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reset lại body
		}

		// Ghi lại response body bằng cách thay thế writer mặc định
		respBody := &bytes.Buffer{}
		writer := &bodyWriter{body: respBody, ResponseWriter: c.Writer}
		c.Writer = writer

		// Gắn id logging vào context
		entry := logrus.WithFields(logrus.Fields{
			"id":     id,
			"source": "service",
		})
		ctx := utils.WithLogger(c.Request.Context(), entry)
		c.Request = c.Request.WithContext(ctx)

		// Tiếp tục xử lý
		c.Next()

		duration := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.RequestURI
		statusCode := c.Writer.Status()
		contentTypeReq := c.ContentType()
		contentTypeResp := c.Writer.Header().Get("Content-Type")
		lang := c.Query("lang")
		accept := c.GetHeader("Accept-Language")
		if lang == "" && accept != "" {
			lang = accept
		}

		// Format body nếu có thể (JSON pretty)
		formattedReq := tryFormatJSON(requestBody)
		var formattedResp string
		if (method == "GET" || method == "") && (hasListQuery(c) || isListRoute(c)) {
			formattedResp = ""
		} else {
			formattedResp = tryFormatJSON(respBody.Bytes())
		}

		// Ghi log
		logger.Printf(`
. ID: %s
. Client IP: %s
. Method: %s
. Path: %s
. Language: %s
. Status code: %d
. Time: %v
. Request body (Content type: %s):
%s
. Response body (Content type: %s):
%s
--------------------------------------------------------------------------
`, id, clientIP, method, path, lang, statusCode, duration, contentTypeReq, formattedReq, contentTypeResp, formattedResp)
	}
}

// bodyWriter để ghi lại response body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // ghi vào bộ nhớ tạm
	return w.ResponseWriter.Write(b)
}

// Format JSON đẹp nếu có thể
func tryFormatJSON(data []byte) string {
	var out bytes.Buffer
	if json.Valid(data) {
		if err := json.Indent(&out, data, "", "  "); err == nil {
			return out.String()
		}
	}
	return string(data)
}

// Lấy thư mục cha từ đường dẫn file
func getDir(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return "."
	}
	return path[:idx]
}

func hasListQuery(c *gin.Context) bool {
	query := c.Request.URL.Query()

	// Nếu có ít nhất một trong các query sau: limit, page, sort, search => coi là danh sách
	return query.Has("limit") || query.Has("page") || query.Has("sort") || query.Has("search")
}

func isListRoute(c *gin.Context) bool {
	p := c.FullPath() // ví dụ: "/api/v1/users" hoặc "/api/v1/users/:id"
	if p == "" {
		// không match route nào (404) hoặc middleware rất sớm
		return false
	}
	// Coi là "list" nếu segment cuối KHÔNG phải tham số (":xxx")
	parts := strings.Split(strings.Trim(p, "/"), "/")
	last := parts[len(parts)-1]
	return !strings.HasPrefix(last, ":")
}
