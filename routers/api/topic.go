package api

import (
	"github.com/gin-gonic/gin"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/pkg/setting"
	"justus/service/statistics_service"
	"justus/service/topic_service"
	"net/http"
)

//获取最热话题
func TopicHotList(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct {
		PageNum int `json:"page_num"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	list, err := topic_service.GetHotList(c.GetString("lange"), request.PageNum, setting.AppSetting.PageSize, "")
	if err != nil {
		appG.Response(e.ERROR, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

//获取最热话题 创建网页拍用 不带#号
func GetWebTopicHotList(c *gin.Context) {
	appG := app.Gin{C: c}

	list, err := topic_service.GetHotListV2(c.GetString("lange"), 1, setting.AppSetting.PageSize, "")
	if err != nil {
		appG.Response(e.ERROR, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

//话题搜索
func TopicSearch(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct {
		KeyWord string `json:"keyword"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	list, err := topic_service.GetSearchList(c.GetString("lange"), 1, setting.AppSetting.PageSize, request.KeyWord)
	if err != nil {
		appG.Response(e.ERROR, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

//热门推荐
func TopicRecommend(c *gin.Context) {
	appG := app.Gin{C: c}
	list, err := topic_service.GetHotList(c.GetString("lange"), 1, 10, "")
	if err != nil {
		appG.Response(e.ERROR, e.ERROR, nil)
		return
	}
	if len(list.List) > 0 {
		go func() { //添加曝光
			var topicIds []int
			for _, v := range list.List {
				topicIds = append(topicIds, v.Id)
			}
			statistics_service.UploadStatisticsByIds(topicIds, statistics_service.FromTopic, statistics_service.Impression)
		}()
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

//获取话题详情
func TopicDetail(c *gin.Context) {
	var err error
	var detail *topic_service.Detail
	appG := app.Gin{C: c}
	type Request struct {
		TopicId   int    `json:"topic_id"`
		TopicName string `json:"topic_name"`
	}
	request := Request{}
	err = c.ShouldBindJSON(&request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	uid := c.GetInt("userId")

	if request.TopicId > 0 {
		detail, err = topic_service.GetDetail(uid, request.TopicId)
	} else {
		detail, err = topic_service.GetDetailByName(uid, request.TopicName)
	}

	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	if detail.Id == 0 {
		appG.Response(http.StatusOK, e.CONTENT_EMPTY, nil)
	} else {
		go func() {
			statistics_service.UploadStatisticsById(request.TopicId, statistics_service.FromTopic, statistics_service.Click)
		}()
		appG.Response(http.StatusOK, e.SUCCESS, detail)
	}
}

//话题收藏
func TopicCollect(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct {
		TopicId int `json:"topic_id"`
		Status  int `json:"status"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	uid := c.GetInt("userId")
	res := topic_service.CollectTopic(uid, request.TopicId, request.Status)
	if res {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	} else {
		appG.Response(e.ERROR, e.ERROR, nil)
	}

}

//话题下的图片
func TopicPictureList(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct {
		TopicId int `json:"topic_id"`
		PageNum int `json:"page_num"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	uid := c.GetInt("userId")
	list, err := topic_service.GetTopicPictureList(request.TopicId, uid, request.PageNum, setting.AppSetting.PageSize)
	if err != nil {
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, list)
}
