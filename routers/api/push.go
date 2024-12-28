package api

import (
	"github.com/gin-gonic/gin"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"net/http"
)

func PushSwitch(c *gin.Context) {
	var err error
	appG := app.Gin{C: c}
	type Request struct{
		PushType int `json:"push_type"`
		Status int `json:"status"`
	}
	request := Request{}
	err = c.ShouldBindJSON(&request)
	if err != nil{
		appG.Response(http.StatusBadRequest, e.ERROR,nil)
		return
	}
	uid := c.GetInt("userId")
	redisKey := ""
	if request.PushType == 1{	//好友推送
		redisKey = rediskey.GetPushMessageFriendSwitchKey(uid)
	}
	if request.PushType == 2{	//热拍推送
		redisKey = rediskey.GetPushMessageHotBeatSwitchKey(uid)
	}
	if request.Status == 1 || request.Status == 0{
		_ = gredis.Set(redisKey, request.Status, 0)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
