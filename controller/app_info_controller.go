package controller

import (
	"net/http"
	"runtime"

	"github.com/elskow/PlagiProof/util"
	"github.com/gin-gonic/gin"
)

type (
	AppInfoController interface {
		HealthCheck(ctx *gin.Context)
		GoVersion(ctx *gin.Context)
	}

	appInfoController struct{}
)

func NewAppInfoController() appInfoController {
	return appInfoController{}
}

func (c *appInfoController) HealthCheck(ctx *gin.Context) {
	res := util.Response{
		Ok:      true,
		Message: "Its Working!",
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *appInfoController) GoVersion(ctx *gin.Context) {
	res := util.Response{
		Ok:      true,
		Message: runtime.Version(),
	}

	ctx.JSON(http.StatusOK, res)
}
