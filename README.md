# 📚 Dự án tự học microservice - 🇬​​🇴​🇱​🇦​🇳​🇬

Go là ngôn ngữ lập trình mã nguồn mở giúp bạn dễ dàng xây dựng phần mềm đơn giản, đáng tin cậy và hiệu quả.

![Gopher image](https://golang.org/doc/gopher/fiveyears.jpg)
_Gopher image by [Renee French][rf], licensed under [Creative Commons 4.0 Attribution license][cc4-by]._

### 🛠️ Cài đặt (Linux/Ubuntu)

- Các bản phân phối nhị phân chính thức có sẵn tại https://go.dev/dl/.

## 🥃 Gin Web Framework

Gin là một framework web được viết bằng Go. Nó có API tương tự Martini nhưng hiệu suất nhanh hơn tới 40 lần nhờ sử dụng httprouter. Nếu bạn cần hiệu năng cao và năng suất tốt, bạn sẽ yêu thích Gin.

<img align="right" width="36%" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">

### Từ khóa liên quan tới những module đã được áp dụng:

1. [Lớp khởi chạy (main.go) 📌](#1-lớp-khởi-chạy-maingo-📌)
2. [Thực thể (Entity) 👤](#2-thực-thể-entity-👤)
3. [Cơ sở dữ liệu (Database) 🛢️](#3-cơ-sở-dữ-liệu-database-🛢️)
4. [ORM 🔄](#4-orm-🔄)
5. [Bộ định tuyến (Router) 📡](#5-bộ-định-tuyến-router-📡)
6. [Phân trang (Pagination) 🔢](#6-phân-trang-pagination-🔢)
7. [Ánh xạ dữ liệu (Mapping) 🔁](#7-ánh-xạ-dữ-liệu-mapping-🔁)
8. [Xác thực & phân quyền 🔐](#8-xác-thực-&-phân-quyền-🔐)
9. [Nhật kí (Logging) 📝](#9-nhật-kí-logging-📝)
10. [Xử lí lỗi toàn cục (Error handler) ⚠️](#10-xử-lí-lỗi-toàn-cục-error-handler-⚠️)
11. [Gỡ lỗi (Debug) 🐞](#11-gỡ-lỗi-debug-🐞)
12. [Validation ✅](#12-validation-✅)
13. [Swagger UI 🍀](#13-swagger-ui-🍀)
14. [gRPC 🔀](#14-grpc-🔀)
15. [Testing 🧪](#15-testing-🧪)
16. [Cache 💾](#16-cache-💾)
17. [Vault 🛡️](#17-vault-🛡️)
18. [Internationalization (I18n) 🌎](#18-internationalization-i18n-🌎)
19. [Cloud service ☁️](#19-cloud-service-☁️)
20. [Deploy & CICD 🚀](#20-deploy--cicd-🚀)

---

## 1. Lớp khởi chạy (main.go) 📌

<!-- Mô tả hoặc ví dụ về lớp khởi chạy -->

```md
- Ví dụ về lớp khởi chạy với Gin framework
- Mục đích của lớp khởi chạy là để khởi tạo các thành phần cần thiết của ứng dụng như môi trường, logger, kết nối cơ sở dữ liệu, và cấu hình các route. Đây là điểm bắt đầu của ứng dụng, nơi mà tất cả các thành phần khác được kết nối với nhau.
- Lớp khởi chạy này sử dụng Gin framework để tạo một HTTP server, kết nối tới cơ sở dữ liệu, và thiết lập các route cho ứng dụng.
- Nó cũng bao gồm việc khởi tạo các biến môi trường, logger, và i18n (internationalization) để hỗ trợ đa ngôn ngữ.
- Cuối cùng, nó chạy server trên cổng mặc định 8080.
```

<details>
<summary>✨ Xem ví dụ đầy đủ</summary>

```go
package main

import (
	"go-demo-gin/docs"
	"go-demo-gin/initializers"
	"go-demo-gin/routes"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// 1. Env
	initializers.LoadEnvVariables()
	if err := initializers.RequireEnv("DB_URL", "SECRET"); err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Missing required environment")
	}

	// 2. Logger
	initializers.InitLogger()

	if err := initializers.LoadI18n(); err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Failed to load i18n")
	}

	// 3. Database
	db, err := initializers.ConnectToDB()
	if err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Fail to connect to database")
	}

	router := gin.Default()

	// Swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Routes (DI)
	routes.SetupRoutes(router, db)

	// (Optional) graceful shutdown: đóng sqlDB khi app dừng
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	router.Run()
}
```

</details>

---

## 2. Thực thể (Entity) 👤

<!-- Mô tả hoặc ví dụ về Entity -->

---

## 3. Cơ sở dữ liệu (Database) 🛢️

<!-- Mô tả hoặc ví dụ về Database -->

---

## 4. ORM 🔄

<!-- Mô tả hoặc ví dụ về ORM -->

---

## 5. Bộ định tuyến (Router) 📡

<!-- Mô tả hoặc ví dụ về Router -->

---

## 6. Phân trang (Pagination) 🔢

<!-- Mô tả hoặc ví dụ về Pagination -->

---

## 7. Ánh xạ dữ liệu (Mapping) 🔁

<!-- Mô tả hoặc ví dụ về Mapping -->

---

## 8. Xác thực & phân quyền 🔐

<!-- Mô tả hoặc ví dụ về xác thực & phân quyền -->

---

## 9. Nhật kí (Logging) 📝

<!-- Mô tả hoặc ví dụ về Logging -->

---

## 10. Xử lí lỗi toàn cục (Error handler) ⚠️

<!-- Mô tả hoặc ví dụ về Error handler -->

---

## 11. Gỡ lỗi (Debug) 🐞

<!-- Mô tả hoặc ví dụ về Debug -->

---

## 12. Validation ✅

<!-- Mô tả hoặc ví dụ về Validation -->

---

## 13. Swagger UI 🍀

<!-- Mô tả hoặc ví dụ về Swagger UI -->

---

## 14. gRPC 🔀

<!-- Mô tả hoặc ví dụ về gRPC -->

---

## 15. Testing 🧪

<!-- Mô tả hoặc ví dụ về Testing -->

---

## 16. Cache 💾

<!-- Mô tả hoặc ví dụ về Cache -->

---

## 17. Vault 🛡️

<!-- Mô tả hoặc ví dụ về Vault -->

---

## 18. Internationalization (I18n) 🌎

<!-- Mô tả hoặc ví dụ về I18n -->

---

## 19. Cloud service ☁️

<!-- Mô tả hoặc ví dụ về Cloud service -->

---

## 20. Deploy & CICD 🚀

<!-- Mô tả hoặc ví dụ về Deploy & CICD -->
