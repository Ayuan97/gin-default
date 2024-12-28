package api

import (
	"github.com/gin-gonic/gin"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service/statistics_service"
	"net/http"
)

//分享上传
func ShareUpload(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct{
		FromType int `json:"from_type"`
		Id int `json:"id"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil{
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS,nil)
		return
	}
	if request.FromType == 1 || request.FromType == 2{	//话题
		_ = statistics_service.UploadStatisticsById(request.Id, request.FromType, statistics_service.Share)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
