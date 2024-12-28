package api

import (
	"github.com/gin-gonic/gin"
	"justus/pkg/app"
	"justus/pkg/e"
	"net/http"
)

func LangeGenerate(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}
