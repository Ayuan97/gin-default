package api

import (
	"github.com/gin-gonic/gin"
	"justus/global"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service"
	"net/http"
)

func Like(c *gin.Context) {
	appG := app.Gin{C: c}
	var param service.PostLike
	userId := c.MustGet("userId").(int)
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("UserCenter BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return

	}
	data, err := service.Like(param, userId)
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

// GetLikeList 获取点赞列表
func GetLikeList(c *gin.Context) {
	appG := app.Gin{C: c}
	var param service.PostLike
	userId := c.MustGet("userId").(int)
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("UserCenter BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return

	}
	//service.LikeList(param, userId)
	data, err := service.LikeList(param, userId)
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
