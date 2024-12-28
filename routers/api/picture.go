package api

import (
	"github.com/gin-gonic/gin"
	"justus/global"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service/picture_service"
	"net/http"
)

func GetPicture(c *gin.Context) {
	appG := app.Gin{C: c}
	var param picture_service.PostPicture
	userId := c.MustGet("userId").(int)
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("UserCenter BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return

	}
	//picture_service.GetPicture(param, userId)
	data, err := picture_service.GetPicture(param, userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
		return
	}
	if data != nil {
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.ERROR, nil)
	return
}

//图片列表
func GetPictureList(c *gin.Context) {
	appG := app.Gin{C: c}
	var param picture_service.PostImgList
	userId := c.MustGet("userId").(int)
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("UserCenter BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return

	}
	//picture_service.GetPictureList(param, userId)
	data, err := picture_service.GetPictureList(param, userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
		return
	}
	if data != nil {
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.ERROR, nil)
	return
}
