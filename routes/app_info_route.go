package routes

import (
	"github.com/elskow/PlagiProof/controller"
	"github.com/gin-gonic/gin"
)

func AppInfoRoute(route *gin.Engine, c controller.AppInfoController) {
	routes := route.Group("/app-info")
	{
		routes.GET("/health-check", c.HealthCheck)
		routes.GET("/go-version", c.GoVersion)
	}
}
