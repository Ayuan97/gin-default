package routers

import (
	"github.com/gin-gonic/gin"
	"justus/middleware/api_require"
	"justus/middleware/bodyLog"
	"justus/middleware/cors"
	"justus/middleware/jwt"
	"justus/middleware/recovers"
	"justus/routers/api"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), recovers.Recover(), bodyLog.GinBodyLogMiddleware())

	apiGroup := r.Group("/api/v2")
	apiGroup.Use(api_require.Common())
	//apiGroup.Use(sign.VerifySignature())
	apiGroup.Use(jwt.JWT())
	{
		apiGroup.POST("/img/list", api.GetPictureList)               //图片列表
		apiGroup.POST("/recommend/get_list", api.RecommendList)      //推荐列表
		apiGroup.POST("/home/index", api.UserCenter)                 //个人主页
		apiGroup.POST("/home/img_list", api.UserCenterImgList)       //个人主页 图片列表
		apiGroup.POST("/topic/hot_list", api.TopicHotList)           //话题热度列表
		apiGroup.POST("/topic/web/hot_list", api.GetWebTopicHotList) //网页拍-话题随机列表
		apiGroup.POST("/topic/search", api.TopicSearch)              //话题热度列表
		apiGroup.POST("/topic/recommend", api.TopicRecommend)        //话题热度列表
		apiGroup.POST("/topic/detail", api.TopicDetail)              //话题详情
		apiGroup.POST("/topic/collect", api.TopicCollect)            //话题收藏
		apiGroup.POST("/topic/picture_list", api.TopicPictureList)   //话题下的图片列表
		apiGroup.POST("/img/like", api.Like)                         //点赞
		apiGroup.POST("/img/like_list", api.GetLikeList)             //点赞列表
		apiGroup.POST("/comment/creat", api.CreateComment)           //评论
		apiGroup.POST("/comment/del", api.DelComment)                //评论删除
		apiGroup.POST("/comment/get_list", api.GetComment)           //评论列表
		apiGroup.POST("/follow/get_list", api.FollowList)            //关注列表
		apiGroup.POST("/user/follow", api.UserFollow)                //用户关注
		apiGroup.POST("/message/list", api.MessageList)              //消息列表
		apiGroup.POST("/message/clear", api.MessageClear)            //清空消息
		apiGroup.POST("/message/delete", api.MessageDelete)          //删除消息
		//apiGroup.POST("/lange/generate", api.LangeGenerate)         //生成语言
		apiGroup.POST("/impression/upload", api.ImpressionUpload) //曝光上传
		apiGroup.POST("/share/upload", api.ShareUpload)           //分享上传
		apiGroup.POST("/recommend/list", api.RecommendList)       //推荐列表
		apiGroup.POST("/picture/get", api.GetPicture)             //获取单张图片信息
		apiGroup.POST("/push/switch", api.PushSwitch)             //推送开关

		//test
		apiGroup.POST("/test", api.Test)
	}
	webapiGroup := r.Group("/api/v2")

	webapiGroup.Use(cors.Cors())
	webapiGroup.Use(api_require.Common())
	webapiGroup.Use(jwt.JWT())
	{
		webapiGroup.GET("/component/get_latest_picture", api.GetLatestPicture)
		webapiGroup.POST("/component/get_latest_picture", api.PostLatestPicture)
	}

	webapi := r.Group("/api/v2")
	webapi.Use(cors.Cors())
	{
		webapi.POST("/web/topic/recommend", api.WebTopicRecommend) //web端话题推荐
	}

	return r
}
