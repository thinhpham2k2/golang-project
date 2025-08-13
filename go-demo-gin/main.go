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
// @description Enter your Bearer token
// @description Example: Bearer 1234567890abcdef
func main() {
	// 1. Env
	initializers.LoadEnvVariables()
	if err := initializers.RequireEnv("DB_URL", "SECRET"); err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Missing required environment")
	}

	// 2. Logger
	initializers.InitLogger()

	// 3. I18n
	if err := initializers.LoadI18n(); err != nil {
		logrus.WithField("source", "system").WithError(err).Fatal("Failed to load i18n")
	}

	// 4. Database
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
