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
