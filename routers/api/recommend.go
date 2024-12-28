package api

import (
	"github.com/gin-gonic/gin"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service"
	"justus/service/recommend_service"
	"net/http"
)

// RecommendList  @yuan 推荐列表
func RecommendList(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := c.MustGet("userId").(int)
	lange := c.MustGet("lange").(string)
	param := service.UserCenter{}
	_ = c.ShouldBindJSON(&param)
	result := recommend_service.RecommendList(userId, param.Page, lange)
	appG.Response(http.StatusOK, e.SUCCESS, result)
}
