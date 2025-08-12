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
