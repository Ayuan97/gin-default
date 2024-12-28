package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"justus/global"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/service/topic_service"
	"net/http"
)

type WebShareLinkData struct {
	TopicName string `json:"topic_name"`
	Uid       int    `json:"uid"`
	Lange     string `json:"lange"`
}

//热门推荐
func WebTopicRecommend(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct {
		Key string `json:"key"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		global.Logger.Error("[WebTopicRecommend] should bind json error:", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	data := gredis.Get(rediskey.GetWebShareLinkKey(request.Key))
	linkData := WebShareLinkData{}
	if data != "" {
		err = json.Unmarshal([]byte(data), &linkData)
		if err != nil {
			global.Logger.Error("[WebTopicRecommend] json 2 unmarshal error:", err)
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
	} else {
		fmt.Println("[WebTopicRecommend] redis get key:", request.Key, "data is empty")
		global.Logger.Error("[WebTopicRecommend] redis get key:", request.Key, "data is empty")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	list, err := topic_service.GetWebHotList(linkData.Lange, 1, 9)
	if err != nil {
		appG.Response(e.ERROR, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, list)
}
