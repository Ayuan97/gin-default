package api

import (
	"github.com/gin-gonic/gin"
	"justus/global"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service"
	"net/http"
)

func CreateComment(c *gin.Context) {
	appG := app.Gin{C: c}
	var param service.PostComment
	userId := c.MustGet("userId").(int)
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("CreateComment BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	data, err := service.CreateComment(param, userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
		return
	}
	if data.ID <= 0 {
		appG.Response(http.StatusOK, e.SUCCESS, make(map[string]interface{}))
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}

func DelComment(c *gin.Context) {
	appG := app.Gin{C: c}
	var param service.PostDelComment
	userId := c.MustGet("userId").(int)
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("DelComment BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	data, err := service.DelComment(param, userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
		return
	}
	if data {
		appG.Response(http.StatusOK, e.SUCCESS, make(map[string]interface{}))
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}

// GetComment 获取评论
func GetComment(c *gin.Context) {
	appG := app.Gin{C: c}
	var param service.PostComment
	userId := c.MustGet("userId").(int)

	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("GetComment BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	data, err := service.GetComment(param.Pid, param.Page, userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}
