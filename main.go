package main

import (
	"os"

	"github.com/elskow/PlagiProof/controller"
	"github.com/elskow/PlagiProof/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	server := gin.Default()
	server.Use(gin.Logger())

	var (
		appInfoController = controller.NewAppInfoController()
	)

	routes.AppInfoRoute(server, &appInfoController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serve := os.Getenv("APP_ENV")
	if serve == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		logrus.Fatalf("error running server: %v", err)
	}
}
