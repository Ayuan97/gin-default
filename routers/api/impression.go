package api

import (
	"github.com/gin-gonic/gin"
	"justus/dao"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service/statistics_service"
	"net/http"
	"strings"
)

//曝光
//topic_list(热门列表页)
//picture_detail(图片详情)
//recommend_list(推荐页)
//follow_list(关注列表)
//topic_detail(话题详情)
func ImpressionUpload(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct{
		FromType int `json:"from_type"`
		Page string `json:"page"`
		TopicIds string `json:"topic_ids"`
		PictureIds string `json:"picture_ids"`
		PictureUIds string `json:"picture_uids"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil{
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS,nil)
		return
	}
	if request.FromType == 1{	//话题
		_ = statistics_service.UploadStatistics(request.TopicIds, "",request.FromType, statistics_service.Impression)
	}
	if request.FromType == 2{	//图片
		_ = statistics_service.UploadStatistics(request.PictureIds,request.PictureUIds, request.FromType, statistics_service.Impression)
	}
	if request.FromType == 3{	//话题图片都有
		_ = statistics_service.UploadStatistics(request.TopicIds, "",1, statistics_service.Impression)
		_ = statistics_service.UploadStatistics(request.PictureIds, request.PictureUIds,2, statistics_service.Impression)
	}
	if request.Page == "recommend_list"{
		uid := c.GetInt("userId")
		if request.PictureIds != ""{
			pictureIds := strings.Split(request.PictureIds,",")
			dao.AddRecommendListReadRecord(uid,pictureIds)
		}

	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
