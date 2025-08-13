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

1. [Lớp khởi chạy (main.go) 📌](#1-lớp-khởi-chạy-maingo-)
2. [Thực thể (Entity) 🎭](#2-thực-thể-entity-)
3. [Cơ sở dữ liệu (Database) 🛢](#3-cơ-sở-dữ-liệu-database-)
4. [ORM 🔄](#4-orm-)
5. [Bộ định tuyến (Router) 📡](#5-bộ-định-tuyến-router-)
6. [Phân trang (Pagination) 🔢](#6-phân-trang-pagination-)
7. [Ánh xạ dữ liệu (Mapping) 🔁](#7-ánh-xạ-dữ-liệu-mapping-)
8. [Xác thực & phân quyền 🔐](#8-xác-thực--phân-quyền-)
9. [Nhật kí (Logging) 📝](#9-nhật-kí-logging-)
10. [Xử lí lỗi toàn cục (Error handler) 🛑](#10-xử-lí-lỗi-toàn-cục-error-handler-)
11. [Gỡ lỗi (Debug) 🐞](#11-gỡ-lỗi-debug-)
12. [Validation ✅](#12-validation-)
13. [Swagger UI 🍀](#13-swagger-ui-)
14. [gRPC 🔀](#14-grpc-)
15. [Testing 🧪](#15-testing-)
16. [Cache 💾](#16-cache-)
17. [Vault 🔰](#17-vault-)
18. [Internationalization (I18n) 🌎](#18-internationalization-i18n-)
19. [Cloud service ⛅](#19-cloud-service-)
20. [Deploy & CICD 🚀](#20-deploy--cicd-)

---

### 1. Lớp khởi chạy (main.go) 📌

<!-- Mô tả hoặc ví dụ về lớp khởi chạy -->

Chạy ứng dụng bằng lệnh:

```
go run main.go
```

Hoặc biên dịch và chạy:

```
go build main.go && ./main
```

Để tự động theo dõi thay đổi và tái biên dịch, bạn có thể sử dụng CompileDaemon:

```
CompileDaemon -command="./main"
```

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
// File: main.go
// Package main là điểm bắt đầu của ứng dụng Go sử dụng Gin framework.
// Nó khởi tạo các thành phần cần thiết như môi trường, logger, kết nối cơ sở dữ liệu, và cấu hình các route.
// Đây là nơi mà ứng dụng được khởi chạy và các thành phần khác được kết nối với nhau.
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

### 2. Thực thể (Entity) 🎭

<!-- Mô tả hoặc ví dụ về Entity -->

```md
- Ví dụ về thực thể User
- Thực thể User đại diện cho người dùng trong hệ thống, bao gồm các trường như Username, Password, Name, Birthday và Role.
- Nó sử dụng GORM để ánh xạ tới bảng người dùng trong cơ sở dữ liệu.
- Các trường được định nghĩa với kiểu dữ liệu phù hợp và các thuộc tính cần thiết để lưu trữ thông tin người dùng.
```

<details>
<summary>✨ Xem ví dụ entity</summary>

```go
// File: models/user.go
// Package models chứa các thực thể của ứng dụng, trong đó có User.
// Thực thể User đại diện cho người dùng trong hệ thống, bao gồm các trường như Username, Password, Name, Birthday và Role.
// Nó sử dụng GORM để ánh xạ tới bảng người dùng trong cơ sở dữ liệu.
// Các trường được định nghĩa với kiểu dữ liệu phù hợp và các thuộc tính cần thiết để lưu trữ thông tin người dùng.
package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Name     sql.NullString
	Birthday *time.Time `gorm:"type:date"`
	Role     Role       `gorm:"type:varchar(20)"`
}
```

</details>

---

### 3. Cơ sở dữ liệu (Database) 🛢

<!-- Mô tả hoặc ví dụ về Database -->

#### Tạo kết nối cơ sở dữ liệu

```md
- Ví dụ về kết nối cơ sở dữ liệu
- Kết nối cơ sở dữ liệu sử dụng GORM, một ORM phổ biến trong Go, để tương tác với cơ sở dữ liệu.
- Hàm ConnectToDB thiết lập kết nối tới cơ sở dữ liệu sử dụng chuỗi kết nối được cung cấp trong biến môi trường DB_URL.
- Nếu kết nối thành công, nó sẽ tự động thực hiện các thao tác cần thiết như tự động tạo bảng dựa trên các thực thể đã định nghĩa.
```

<details>
<summary>✨ Xem ví dụ khởi tạo kết nối Database</summary>

```go
// File: initializers/db.go
// Package initializers chứa các hàm khởi tạo cho ứng dụng, bao gồm kết nối cơ sở dữ liệu.
// Hàm ConnectToDB thiết lập kết nối tới cơ sở dữ liệu sử dụng chuỗi kết nối được cung cấp trong biến môi trường DB_URL.
// Nếu kết nối thành công, nó sẽ tự động thực hiện các thao tác cần thiết như tự động tạo bảng dựa trên các thực thể đã định nghĩa.
// Kết nối cơ sở dữ liệu sử dụng GORM, một ORM phổ biến trong Go, để tương tác với cơ sở dữ liệu.
package initializers

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, errors.New("DB_URL is empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Pooling (tuỳ chỉnh qua env nếu muốn)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Check kết nối thực
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	logrus.WithField("source", "system").Info("Connected to database")
	return db, nil
}
```

</details>

#### Database migration

Chạy lệnh migrate để tự động tạo bảng:

```
go run migrate/migrate.go
```

```md
- Ví dụ về tự động tạo bảng
- Hàm AutoMigrate được sử dụng để tự động tạo bảng trong cơ sở dữ liệu dựa trên các thực thể đã định nghĩa.
- Nó sẽ kiểm tra và tạo bảng cho thực thể User nếu nó chưa tồn tại.
- Điều này giúp đảm bảo rằng cơ sở dữ liệu luôn được cập nhật với các thay đổi trong mô hình dữ liệu mà không cần phải viết các câu lệnh SQL thủ công.
```

<details>
<summary>✨ Xem ví dụ tự động tạo bảng (Database first)</summary>

```go
// File: migrate/migrate.go
// Package migrate chứa các hàm để thực hiện các thao tác di chuyển cơ sở dữ liệu, bao gồm tự động tạo bảng dựa trên các thực thể đã định nghĩa.
// Hàm main là điểm bắt đầu của ứng dụng, nơi mà các thành phần cần thiết được khởi tạo và kết nối tới cơ sở dữ liệu.
// Nó sử dụng GORM để tương tác với cơ sở dữ liệu và tự động tạo bảng cho thực thể User.
package main

import (
	"go-demo-gin/initializers"
	"go-demo-gin/models"

	"github.com/sirupsen/logrus"
)

func main() {
	// 1. Env
	initializers.LoadEnvVariables()
	if err := initializers.RequireEnv("DB_URL"); err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Missing required environment")
	}

	// 2. Database
	db, err := initializers.ConnectToDB()
	if err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Fail to connect to database")
	}

	db.AutoMigrate(&models.User{})
}
```

</details>

---

### 4. ORM 🔄

<!-- Mô tả hoặc ví dụ về ORM -->

```md
- Ví dụ về ORM với GORM
- GORM là một ORM (Object Relational Mapping) phổ biến trong Go, giúp tương tác với cơ sở dữ liệu một cách dễ dàng và hiệu quả.
- Nó cung cấp các phương thức để thực hiện các thao tác CRUD (Create, Read, Update, Delete) trên các thực thể đã định nghĩa.
- GORM hỗ trợ nhiều loại cơ sở dữ liệu khác nhau và cho phép ánh xạ các trường trong thực thể tới các cột trong bảng cơ sở dữ liệu.
```

<details>
<summary>✨ Xem ví dụ về ORM</summary>

```go
// File: repo/user.go
// Package repo chứa các kho lưu trữ dữ liệu, trong đó có UserRepo.
// GormUserRepo là một triển khai của UserRepo sử dụng GORM để tương tác với cơ sở dữ liệu.
// Nó cung cấp các phương thức để thực hiện các thao tác CRUD (Create, Read, Update, Delete) trên thực thể User.
// Các phương thức này sử dụng GORM để thực hiện các truy vấn tới cơ sở dữ liệu một cách dễ dàng và hiệu quả.
package repo

import (
	"context"

	"go-demo-gin/models"
	"go-demo-gin/pkg"
	"go-demo-gin/utils"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, tx *gorm.DB, u *models.User) error
	FindByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User, error)
	Update(ctx context.Context, tx *gorm.DB, u *models.User) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	List(ctx context.Context, tx *gorm.DB, pag *pkg.Pagination, search string) ([]models.User, int64, error)
	FindByUsername(ctx context.Context, tx *gorm.DB, username string) (*models.User, error)
}

type GormUserRepo struct{ db *gorm.DB }

func NewGormUserRepo(db *gorm.DB) *GormUserRepo { return &GormUserRepo{db: db} }

func (r *GormUserRepo) use(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *GormUserRepo) Create(ctx context.Context, tx *gorm.DB, u *models.User) error {
	return r.use(tx).WithContext(ctx).Create(u).Error
}

func (r *GormUserRepo) FindByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User, error) {
	var u models.User
	if err := r.use(tx).WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepo) Update(ctx context.Context, tx *gorm.DB, u *models.User) error {
	return r.use(tx).WithContext(ctx).Updates(u).Error
}

func (r *GormUserRepo) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return r.use(tx).WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *GormUserRepo) List(ctx context.Context, tx *gorm.DB, pag *pkg.Pagination, search string) ([]models.User, int64, error) {
	q := r.use(tx).WithContext(ctx).Model(&models.User{})
	if search != "" {
		q = q.Where("name ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%") // Postgres: ILIKE
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []models.User
	if err := q.Scopes(utils.Paginate(pag, q)).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormUserRepo) FindByUsername(ctx context.Context, tx *gorm.DB, username string) (*models.User, error) {
	var u models.User
	if err := r.use(tx).WithContext(ctx).
		Where("username = ?", username).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
```

</details>

---

### 5. Bộ định tuyến (Router) 📡

<!-- Mô tả hoặc ví dụ về Router -->

```md
- Ví dụ về bộ định tuyến với Gin framework
- Bộ định tuyến sử dụng Gin framework để định nghĩa các route và ánh xạ chúng tới các controller.
- Nó sử dụng Dependency Injection (DI) để tạo các dịch vụ và controller cần thiết cho ứng dụng.
- Các route được phân chia theo phiên bản API (v1) và các nhóm chức năng (users, authen).
- Mỗi route được bảo vệ bởi các middleware để xác thực và phân quyền người dùng.
- Các route này cho phép thực hiện các thao tác CRUD (Create, Read, Update, Delete) trên thực thể User và các chức năng xác thực người dùng.
```

<details>
<summary>✨ Xem ví dụ về router</summary>

```go
// File: routes/routes.go
// Package routes chứa các hàm để thiết lập các route cho ứng dụng.
// Nó sử dụng Gin framework để định nghĩa các route và ánh xạ chúng tới các controller.
// Hàm SetupRoutes là điểm bắt đầu để cấu hình các route, nơi mà các middleware và các route được thiết lập.
// Nó cũng sử dụng Dependency Injection (DI) để tạo các dịch vụ và controller cần thiết cho ứng dụng.
// Các route được phân chia theo phiên bản API (v1) và các nhóm chức năng (users, authen).
// Mỗi route được bảo vệ bởi các middleware để xác thực và phân quyền người dùng.
// Các route này cho phép thực hiện các thao tác CRUD (Create, Read, Update, Delete) trên thực thể User và các chức năng xác thực người dùng.
// Nó cũng sử dụng Swagger để tạo tài liệu API tự động.
// Gắn Access log filter, Error handler, I18n middleware để xử lý các yêu cầu HTTP một cách hiệu quả và dễ dàng.
// Cuối cùng, nó chạy server trên cổng mặc định 8080.
// Đây là nơi mà ứng dụng được cấu hình để xử lý các yêu cầu HTTP từ người dùng.
package routes

import (
	"go-demo-gin/controllers"
	"go-demo-gin/middlewares"
	"go-demo-gin/models"
	"go-demo-gin/services"
	"go-demo-gin/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	// Use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Gắn Access log filter
	r.Use(middlewares.AccessLogger())

	// Gắn middleware error handler
	r.Use(middlewares.ErrorHandler())

	// Gắn middleware i18n
	r.Use(middlewares.I18nMiddleware())

	ADMIN := models.RoleAdmin
	STAFF := models.RoleStaff
	CUSTOMER := models.RoleCustomer
	RequireRoles := middlewares.AuthenticationFilter(db)

	// Dependency Injection (DI) - constructor injection
	// Create a validator
	v := utils.NewValidator(db)

	// Create services and controllers
	// User service and controller
	userSvc := services.NewUserService(db)
	uc := controllers.NewUserController(v, userSvc)

	// Authen service and controller
	// Read JWT secret from environment variable
	cfg := services.AuthConfig{
		JWTKey:    []byte(os.Getenv("SECRET")),
		Issuer:    "go-demo-gin",
		AccessTTL: time.Hour * 24 * 30,
	}
	authenSvc := services.NewAuthenService(db, cfg)
	ac := controllers.NewAuthenController(authenSvc)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.POST("", RequireRoles(ADMIN, STAFF), uc.UsersCreate)
				users.GET("", RequireRoles(ADMIN, STAFF, CUSTOMER), uc.UsersIndex)
				users.GET("/:id", RequireRoles(ADMIN, STAFF, CUSTOMER), uc.UsersShow)
				users.PUT("/:id", RequireRoles(ADMIN, STAFF, CUSTOMER), uc.UsersUpdate)
				users.DELETE("/:id", RequireRoles(ADMIN, STAFF), uc.UsersDelete)
			}
			authen := v1.Group("/authen")
			{
				authen.POST("/login", ac.Login)
			}
		}
	}
}
```

</details>

---

### 6. Phân trang (Pagination) 🔢

<!-- Mô tả hoặc ví dụ về Pagination -->

```md
- Ví dụ về phân trang
- Phân trang là một kỹ thuật để chia nhỏ dữ liệu thành các trang, giúp quản lý và hiển thị dữ liệu hiệu quả hơn.
- Trong Go, phân trang có thể được thực hiện bằng cách sử dụng một cấu trúc dữ liệu để quản lý các thông tin như giới hạn (limit), trang (page), sắp xếp (sort), tổng số hàng (total_rows), tổng số trang (total_pages) và kết quả (result).
- Cấu trúc dữ liệu này cung cấp các phương thức để tính toán offset, giới hạn, trang và sắp xếp.
- Phân trang cũng hỗ trợ việc trả về kết quả dưới dạng một mảng các đối tượng, cho phép dễ dàng hiển thị kết quả trong các API.
```

<details>
<summary>✨ Xem ví dụ về phân trang</summary>

```go
// File: pkg/pagination.go
// Package pkg chứa các tiện ích chung cho ứng dụng, trong đó có Pagination.
// Pagination là một cấu trúc dữ liệu để quản lý phân trang trong các truy vấn tới cơ sở dữ liệu.
// Nó bao gồm các trường như Limit, Page, Sort, TotalRows, TotalPages và Result.
// Pagination cung cấp các phương thức để tính toán offset, giới hạn, trang, và sắp xếp.
// Các phương thức này giúp dễ dàng quản lý phân trang trong các truy vấn tới cơ sở dữ liệu.
// Pagination cũng hỗ trợ việc trả về kết quả dưới dạng một mảng các đối tượng, cho phép dễ dàng hiển thị kết quả trong các API.
package pkg

type Pagination struct {
	Limit      int    `json:"limit,omitempty" form:"limit"`
	Page       int    `json:"page,omitempty" form:"page"`
	Sort       string `json:"sort,omitempty" form:"sort"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Result     any    `json:"result" swaggertype:"array,object"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

// File: repo/user.go
// Ví dụ áp dụng phân trang trong kho lưu trữ người dùng
// GormUserRepo cũng hỗ trợ phân trang trong phương thức List, cho phép lấy danh sách người dùng với phân trang.
// Phương thức List nhận vào một đối tượng Pagination và trả về danh sách người dùng cùng với tổng số lượng người dùng.
// Nó sử dụng phương thức Scopes của GORM để áp dụng phân trang và sắp xếp cho truy vấn.
// Phương thức này cũng hỗ trợ tìm kiếm người dùng theo tên hoặc tên đăng nhập thông qua tham số search.
// Nếu có lỗi xảy ra trong quá trình truy vấn, nó sẽ trả về lỗi tương ứng.
func (r *GormUserRepo) List(ctx context.Context, tx *gorm.DB, pag *pkg.Pagination, search string) ([]models.User, int64, error) {
	q := r.use(tx).WithContext(ctx).Model(&models.User{})
	if search != "" {
		q = q.Where("name ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%") // Postgres: ILIKE
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []models.User
	if err := q.Scopes(utils.Paginate(pag, q)).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
```

</details>

---

### 7. Ánh xạ dữ liệu (Mapping) 🔁

<!-- Mô tả hoặc ví dụ về Mapping -->

```md- Ví dụ về ánh xạ dữ liệu
- Ánh xạ dữ liệu là quá trình chuyển đổi giữa các đối tượng trong ứng dụng và các thực thể trong cơ sở dữ liệu.
- Trong Go, ánh xạ dữ liệu có thể được thực hiện bằng cách sử dụng các thư viện như `copier` để sao chép dữ liệu giữa các cấu trúc khác nhau.
- Quá trình này giúp dễ dàng chuyển đổi giữa các đối tượng dữ liệu đầu vào (DTO) và các thực thể trong cơ sở dữ liệu.
```

<details>
<summary>✨ Xem ví dụ về mapping</summary>

```go
// File: models/user_create.go
// Package models chứa các thực thể của ứng dụng, trong đó có UserCreate.
// Thực thể UserCreate đại diện cho dữ liệu đầu vào khi tạo người dùng mới.
// Nó bao gồm các trường như Username, Password, Name, Birthday và Role.
// Thực thể này cho phép tạo người dùng mới với các thông tin cần thiết.
// Nó cũng sử dụng các phương thức để xử lý password và ngày sinh, bao gồm việc hash password và chuyển đổi ngày sinh từ chuỗi sang kiểu thời gian.
// Các trường được định nghĩa với kiểu dữ liệu phù hợp và các thuộc tính cần thiết để lưu trữ thông tin người dùng.
// Thực thể này cũng sử dụng các phương thức để xử lý password và ngày sinh, bao gồm việc hash password và chuyển đổi ngày sinh từ chuỗi sang kiểu thời gian.
package user

import (
	"go-demo-gin/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserCreate struct {
	Username string `json:"username" validate:"required,username,duplicateUsername"`
	Pass     string `json:"password" validate:"required,password,hashed" default:"12345678"`
	Name     string `json:"full_name"`
	Role     string `json:"role" validate:"required,role" default:"customer"`
	Date     string `json:"birthday" validate:"birthday" default:"2006-01-02"`
}

func (u *UserCreate) Password() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	return string(hash)
}

func (u *UserCreate) Birthday() *time.Time {
	birthday, _ := time.Parse("2006-01-02", u.Date)
	return &birthday
}

// File: models/user_update.go
// Package models chứa các thực thể của ứng dụng, trong đó có UserUpdate.
// Thực thể UserUpdate đại diện cho dữ liệu đầu vào khi cập nhật thông tin người dùng.
// Nó bao gồm các trường như Pass, Name, Role và Date.
// Thực thể này cho phép cập nhật thông tin người dùng mà không cần phải cung cấp tất cả các trường.
// Nếu trường Pass không được cung cấp, nó sẽ giữ nguyên password cũ.
// Trường Date được sử dụng để lưu trữ ngày sinh của người dùng, và nó sẽ được chuyển đổi thành kiểu thời gian (time.Time) khi cần thiết.
// Nó cũng sử dụng các phương thức để xử lý password và ngày sinh, bao gồm việc hash password và chuyển đổi ngày sinh từ chuỗi sang kiểu thời gian.
package user

import (
	"go-demo-gin/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdate struct {
	Pass string `json:"password" validate:"omitempty,password,hashed" default:"12345678"`
	Name string `json:"full_name"`
	Role string `json:"role" validate:"required,role" default:"customer"`
	Date string `json:"birthday" validate:"birthday" default:"2006-01-02"`
}

func (u *UserUpdate) Password() *string {
	if u.Pass == "" {
		return nil // Không hash, giữ nguyên password cũ
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	result := string(hash)
	return &result
}

func (u *UserUpdate) Birthday() *time.Time {
	birthday, _ := time.Parse("2006-01-02", u.Date)
	return &birthday
}

// File: models/user.go
// Package models chứa các thực thể của ứng dụng, trong đó có User.
// Thực thể đích trong mapping object

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Name     sql.NullString
	Birthday *time.Time `gorm:"type:date"`
	Role     Role       `gorm:"type:varchar(20)"`
}

// Mapper
// File: services/user_service.go
var user models.User
copier.Copy(&user, &in)
```

</details>

---

### 8. Xác thực & phân quyền 🔐

<!-- Mô tả hoặc ví dụ về xác thực & phân quyền -->

```md
- Ví dụ về xác thực và phân quyền
- Xác thực và phân quyền là hai khía cạnh quan trọng trong bảo mật ứng dụng.
- Xác thực là quá trình xác minh danh tính của người dùng, thường thông qua việc sử dụng token JWT (JSON Web Token).
- Trong Go, xác thực và phân quyền có thể được thực hiện bằng cách sử dụng middleware để kiểm tra token JWT trong header của yêu cầu HTTP.
- Middleware này sẽ xác minh token, lấy thông tin người dùng từ cơ sở dữ liệu và kiểm tra xem người dùng có quyền truy cập vào route hiện tại hay không.
- Nếu người dùng được xác thực và có quyền truy cập, middleware sẽ lưu thông tin người dùng vào context và tiếp tục xử lý yêu cầu.
- Đây là một phần quan trọng trong việc bảo mật ứng dụng và đảm bảo rằng chỉ những người dùng có quyền mới có thể truy cập vào các route nhất định.
```

<details>
<summary>✨ Xem ví dụ về xác thực & phân quyền</summary>

```go
// File: middlewares/authenticationFilter.go
// Package middlewares chứa các middleware cho ứng dụng, trong đó có AuthenticationFilter.
// AuthenticationFilter là một middleware để xác thực và phân quyền người dùng.
// Nó kiểm tra header Authorization để xác thực token JWT và phân quyền người dùng dựa trên vai trò (Role).
// Middleware này sẽ lấy thông tin người dùng từ cơ sở dữ liệu dựa trên token JWT và kiểm tra xem người dùng có quyền truy cập vào route hiện tại hay không.
// Nếu người dùng không được xác thực hoặc không có quyền truy cập, middleware sẽ trả về lỗi 401 Unauthorized hoặc 403 Forbidden.
// Nếu người dùng được xác thực và có quyền truy cập, nó sẽ lưu thông tin người dùng vào context và tiếp tục xử lý yêu cầu.
// Middleware này sử dụng GORM để truy vấn cơ sở dữ liệu và jwt-go để xác thực token JWT.
// Nó cũng sử dụng các hàm tiện ích từ package utils để tải biến môi trường và thông điệp i18n.
// Đây là một phần quan trọng trong việc bảo mật ứng dụng và đảm bảo rằng chỉ những người dùng có quyền mới có thể truy cập vào các route nhất định.
// Middleware này cũng hỗ trợ đa ngôn ngữ thông qua việc sử dụng i18n để trả về thông điệp lỗi phù hợp với ngôn ngữ của người dùng.
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

func AuthenticationFilter(db *gorm.DB) func(allowedRoles ...models.Role) gin.HandlerFunc {
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

// File: routes/routes.go
// Package routes chứa các hàm để thiết lập các route cho ứng dụng.
// Nó sử dụng Gin framework để định nghĩa các route và ánh xạ chúng tới các controller.
// Mỗi route được bảo vệ bởi các middleware để xác thực và phân quyền người dùng.
ADMIN := models.RoleAdmin
STAFF := models.RoleStaff
CUSTOMER := models.RoleCustomer
RequireRoles := middlewares.AuthenticationFilter(db)
api := r.Group("/api")
{
	v1 := api.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", RequireRoles(ADMIN, STAFF), uc.UsersCreate)
			users.GET("", RequireRoles(ADMIN, STAFF, CUSTOMER), uc.UsersIndex)
			users.GET("/:id", RequireRoles(ADMIN, STAFF, CUSTOMER), uc.UsersShow)
			users.PUT("/:id", RequireRoles(ADMIN, STAFF, CUSTOMER), uc.UsersUpdate)
			users.DELETE("/:id", RequireRoles(ADMIN, STAFF), uc.UsersDelete)
		}
		authen := v1.Group("/authen")
		{
			authen.POST("/login", ac.Login)
		}
	}
}
```

</details>

---

### 9. Nhật kí (Logging) 📝

<!-- Mô tả hoặc ví dụ về Logging -->

#### Ghi chép nhật kí Access log (log)

```md
- Ví dụ về ghi chép nhật kí Access log
- Ghi chép nhật kí là quá trình ghi lại các thông tin quan trọng về các yêu cầu HTTP và phản hồi của ứng dụng.
- Trong Go, việc ghi chép nhật kí có thể được thực hiện bằng cách sử dụng middleware để ghi lại các thông tin như ID, IP client, phương thức HTTP, đường dẫn, ngôn ngữ, mã trạng thái, thời gian xử lý, body yêu cầu và body phản hồi.
- Middleware này sẽ ghi lại các thông tin này vào một file log cụ thể.
- Nếu không có ID trong header, middleware sẽ tự động sinh một ID dựa trên thời gian hiện tại.
- Middleware này cũng hỗ trợ định dạng body yêu cầu và phản hồi dưới dạng JSON đẹp (pretty JSON).
- Nó cũng kiểm tra xem yêu cầu có phải là danh sách hay không dựa trên các tham số truy vấn như limit, page, sort, search.
- Nếu là danh sách, nó sẽ không ghi lại body phản hồi để tránh ghi lại quá nhiều dữ liệu.
- Việc ghi chép nhật kí giúp theo dõi và phân tích các yêu cầu HTTP, phát hiện lỗi và cải thiện hiệu suất của ứng dụng.
```

<details>
<summary>✨ Xem ví dụ về Access log</summary>

```go
// File: middlewares/accessLog.go
// Package middlewares chứa các middleware cho ứng dụng, trong đó có AccessLogger.
// AccessLogger là một middleware để ghi lại nhật kí truy cập (access log) cho các yêu cầu HTTP.
// Nó ghi lại các thông tin như ID, IP client, phương thức HTTP,
// đường dẫn, ngôn ngữ, mã trạng thái, thời gian xử lý, body yêu cầu và body phản hồi.
// Middleware này sẽ ghi lại các thông tin này vào một file log cụ thể.
// Nếu không có ID trong header, nó sẽ tự động sinh một ID dựa trên thời gian hiện tại.
// Middleware này cũng hỗ trợ định dạng body yêu cầu và phản hồi dưới dạng JSON đẹp (pretty JSON).
// Nó cũng kiểm tra xem yêu cầu có phải là danh sách hay không dựa trên các tham số truy vấn như limit, page, sort, search.
// Nếu là danh sách, nó sẽ không ghi lại body phản hồi để tránh ghi lại quá nhiều dữ liệu.
package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		c.Set("id", id)

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
		if (method == "GET" || method == "") && hasListQuery(c) {
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
```

</details>

#### Ghi chép nhật kí App log (logrus)

```md
- Ví dụ về ghi chép nhật kí App log
- Ghi chép nhật kí ứng dụng là quá trình ghi lại các thông tin quan trọng về hoạt động của ứng dụng, bao gồm các thông tin như ID, user ID, và mức độ ghi chép.
- Trong Go, việc ghi chép nhật kí ứng dụng có thể được thực hiện bằng cách sử dụng thư viện logrus để ghi lại các thông tin này vào một file log cụ thể.
- Thư viện logrus hỗ trợ các tính năng như xoay vòng file log, định dạng JSON hoặc văn bản, và cấu hình mức độ ghi chép.
- Việc ghi chép nhật kí ứng dụng giúp theo dõi và phân tích các hoạt động của ứng dụng, phát hiện lỗi và cải thiện hiệu suất của ứng dụng.
- Hàm InitLogger trong package initializers là một ví dụ về cách khởi tạo cấu hình ghi chép nhật kí cho ứng dụng.
- Hàm này sẽ tạo thư mục chứa file log nếu chưa tồn tại, và cấu hình rotator để xoay vòng file log khi nó đạt đến kích thước tối đa.
```

<details>
<summary>✨ Xem ví dụ về App log</summary>

```go
// File: initializers/logger.go
// Package initializers chứa các hàm khởi tạo cho ứng dụng, trong đó có InitLogger.
// InitLogger là một hàm để khởi tạo cấu hình ghi chép nhật kí cho ứng dụng.
// Nó sử dụng thư viện logrus để ghi chép nhật kí với các tính năng như xoay vòng file log, định dạng JSON hoặc văn bản, và cấu hình mức độ ghi chép.
// Hàm này sẽ tạo thư mục chứa file log nếu chưa tồn tại, và cấu hình rotator để xoay vòng file log khi nó đạt đến kích thước tối đa.
// Nó cũng cấu hình định dạng ghi chép nhật kí dựa trên biến môi trường LOG_FORMAT, với các tùy chọn là "text" hoặc "json".
// Mức độ ghi chép nhật kí cũng được cấu hình dựa trên biến môi trường LOG_LEVEL, với các tùy chọn là "debug", "info", "warn", "error".
// Cuối cùng, nó sẽ ghi chép nhật kí vào cả stdout và file log, đảm bảo rằng các thông tin quan trọng được ghi lại một cách hiệu quả và dễ dàng truy cập.
// Hàm này cũng hỗ trợ xoay vòng file log để tránh việc file log quá lớn, và nén các file log cũ để tiết kiệm không gian lưu tr
package initializers

import (
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	logPath := getEnv("LOG_FILE", "log/app.log")
	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)

	rotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    50, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Nếu mở file/rotator lỗi thì vẫn có stdout
	log.SetOutput(io.MultiWriter(os.Stdout, rotator))

	// Dev: text; Prod: JSON (tuỳ env)
	if getEnv("LOG_FORMAT", "text") == "json" {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:    true,
			DisableColors:    true, // vì có ghi ra file
			QuoteEmptyFields: true, // đảm bảo "source" có ngoặc kép
			ForceQuote:       true,
			TimestampFormat:  time.RFC3339Nano,
		})
	}

	// Level qua env: debug|info|warn|error
	switch getEnv("LOG_LEVEL", "info") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// (Tuỳ) in caller cho debug, tốn chút overhead
	// log.SetReportCaller(true)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// File: utils/util.go
// Package utils chứa các hàm tiện ích chung cho ứng dụng, trong đó có Log.
// Hàm Log là một hàm tiện ích để ghi chép nhật kí với các thông tin như request ID, user ID và mức độ ghi chép.
// Hàm này sử dụng thư viện logrus để ghi chép nhật kí với các mức độ khác nhau như Debug, Info, Warn, Error, Fatal, Trace.
// Nó lấy request ID từ context của Gin, và nếu không có thì gán là "unknown".
// Nó cũng lấy thông tin người dùng từ context, và nếu không có thì gán user mặc định (có thể là một struct rỗng hoặc giá trị mặc định của bạn).
func Log(c *gin.Context, level log.Level, message string) {
	// Lấy request id cho logging
	id, exists := c.Get("id")
	if !exists {
		id = "unknown"
	}

	// Lấy thông tin user cho logging
	val, exists := c.Get("user")  // exists: key có tồn tại không
	user, ok := val.(models.User) // ok: ép kiểu có thành công không
	if !exists || !ok {
		// Gán user mặc định (nếu models.User là struct)
		user = models.User{} // hoặc giá trị mặc định của bạn
	}

	entry := log.WithFields(log.Fields{
		"id":      id,
		"user_id": user.ID,
		"source":  "service",
	})

	switch level {
	case log.DebugLevel:
		entry.Debug(message)
	case log.InfoLevel:
		entry.Info(message)
	case log.WarnLevel:
		entry.Warn(message)
	case log.ErrorLevel:
		entry.Error(message)
	case log.FatalLevel:
		entry.Fatal(message)
	case log.TraceLevel:
		entry.Trace(message)
	default:
		entry.Print(message)
	}
}
```

</details>

---

### 10. Xử lí lỗi toàn cục (Error handler) 🛑

<!-- Mô tả hoặc ví dụ về Error handler -->

```md
- Ví dụ về xử lí lỗi toàn cục
- Xử lí lỗi toàn cục là quá trình xử lý các lỗi xảy ra trong ứng dụng một cách nhất quán và hiệu quả.
- Trong Go, việc xử lí lỗi toàn cục có thể được thực hiện bằng cách sử dụng middleware để kiểm tra các lỗi trong context sau khi xử lý yêu cầu.
- Middleware này sẽ kiểm tra xem có lỗi nào trong context hay không, nếu có, nó sẽ trả về phản hồi JSON với mã trạng thái HTTP tương ứng.
- Nếu lỗi là một HTTPError, nó sẽ trả về mã trạng thái và thông điệp đã được định nghĩa trong HTTPError.
- Nếu lỗi là một lỗi thường, nó sẽ trả về mã trạng thái 500 (Internal Server Error) cùng với thông điệp lỗi.
- Middleware này giúp đảm bảo rằng tất cả các lỗi trong ứng dụng đều được xử lý một cách nhất quán và trả về phản hồi phù hợp cho người dùng.
- Việc xử lí lỗi toàn cục giúp cải thiện trải nghiệm người dùng và giảm thiểu các lỗi không mong muốn trong ứng dụng.
```

<details>
<summary>✨ Xem ví dụ về xử lí lỗi toàn cục</summary>

```go
// File: middlewares/errorHandler.go
// Package middlewares chứa các middleware cho ứng dụng, trong đó có ErrorHandler.
// ErrorHandler là một middleware để xử lý lỗi toàn cục trong ứng dụng.
// Nó sẽ kiểm tra các lỗi trong context sau khi xử lý yêu cầu và trả về phản hồi JSON với mã trạng thái HTTP tương ứng.
// Nếu lỗi là một HTTPError, nó sẽ trả về mã trạng thái và thông điệp đã được định nghĩa trong HTTPError.
// Nếu lỗi là một lỗi thường, nó sẽ trả về mã trạng thái 500 (Internal Server Error) cùng với thông điệp lỗi.
// Middleware này giúp đảm bảo rằng tất cả các lỗi trong ứng dụng đều được xử lý một cách nhất quán và trả về phản hồi phù hợp cho người dùng.
// Nó cũng sử dụng package errorResponse để định nghĩa các lỗi HTTP cụ thể, giúp dễ dàng quản lý và xử lý các lỗi trong ứng dụng.
// Middleware này cũng hỗ trợ đa ngôn ngữ thông qua việc sử dụng i18n để trả về thông điệp lỗi phù hợp với ngôn ngữ của người dùng.
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
```

</details>

---

### 11. Gỡ lỗi (Debug) 🐞

<!-- Mô tả hoặc ví dụ về Debug -->

```md

```

<details>
<summary>✨ Xem ví dụ về debug</summary>

```go

```

</details>

---

### 12. Validation ✅

<!-- Mô tả hoặc ví dụ về Validation -->

```md

```

<details>
<summary>✨ Xem ví dụ về validation</summary>

```go

```

</details>

---

### 13. Swagger UI 🍀

<!-- Mô tả hoặc ví dụ về Swagger UI -->

```md

```

<details>
<summary>✨ Xem ví dụ về swagger ui</summary>

```go

```

</details>

---

### 14. gRPC 🔀

<!-- Mô tả hoặc ví dụ về gRPC -->

```md

```

<details>
<summary>✨ Xem ví dụ về gRPC</summary>

```go

```

</details>

---

### 15. Testing 🧪

<!-- Mô tả hoặc ví dụ về Testing -->

```md

```

<details>
<summary>✨ Xem ví dụ về testing</summary>

```go

```

</details>

---

### 16. Cache 💾

<!-- Mô tả hoặc ví dụ về Cache -->

```md

```

<details>
<summary>✨ Xem ví dụ về cache</summary>

```go

```

</details>

---

### 17. Vault 🔰

<!-- Mô tả hoặc ví dụ về Vault -->

```md

```

<details>
<summary>✨ Xem ví dụ về vault</summary>

```go

```

</details>

---

### 18. Internationalization (I18n) 🌎

<!-- Mô tả hoặc ví dụ về I18n -->

```md

```

<details>
<summary>✨ Xem ví dụ về i18n</summary>

```go

```

</details>

---

### 19. Cloud service ⛅

<!-- Mô tả hoặc ví dụ về Cloud service -->

```md

```

<details>
<summary>✨ Xem ví dụ về các cloud service</summary>

```go

```

</details>

---

### 20. Deploy & CICD 🚀

<!-- Mô tả hoặc ví dụ về Deploy & CICD -->

```md

```

<details>
<summary>✨ Xem ví dụ về deploy & CICD</summary>

```go

```

</details>
