package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"justus/dao"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/pkg/setting"
	"justus/service/message_service"
	"justus/service/push_service"
	"justus/service/statistics_service"
	"strconv"
	"time"
)

// PostLike 点赞入参
type PostLike struct {
	Pid    int `json:"pid"`
	Status int `json:"status"`
	Page   int `json:"page"`
}

func Like(pararm PostLike, uid int) (bool, error) {
	key := rediskey.GetPictureLikeMapKey(pararm.Pid)
	if pararm.Status == 1 {
		//点赞
		var z redis.Z
		z.Score = float64(time.Now().Unix())
		z.Member = strconv.Itoa(uid)
		_, err := gredis.Zadd(key, z)
		if err != nil {
			fmt.Println("err", err)
			return false, err
		}

		go func() { //协程 添加消息
			message_service.AddPictureLikeMessage(uid, pararm.Pid)
			statistics_service.UploadStatisticsById(pararm.Pid, statistics_service.FromPicture, statistics_service.Like)
			push_service.SendPushInList(push_service.PushSendTypePictureLike, uid, push_service.PushMessageParam{PictureId: pararm.Pid})
		}()
		return true, nil
	} else {
		//取消点赞
		_, err := gredis.Zrem(key, strconv.Itoa(uid))
		if err != nil {
			return false, err
		}

		go func() { //协程 添加消息
			statistics_service.UploadStatisticsDecById(pararm.Pid, statistics_service.FromPicture, statistics_service.Like)

		}()
		return true, nil
	}

}

// GetLikeList 获取点赞列表
func LikeList(p PostLike, uid int) (map[string][]interface{}, error) {
	var offset int
	if p.Page > 0 {
		offset = (p.Page - 1) * setting.AppSetting.PageSize
	} else {
		offset = 0
	}
	key := rediskey.GetPictureLikeMapKey(p.Pid)
	val, _ := gredis.Zrange(key, int64(offset), int64(setting.AppSetting.PageSize))
	likeList := map[string][]interface{}{}
	likeList["list"] = []interface{}{}
	if len(val) > 0 {
		var uiDs []int
		for _, v := range val {
			uid, _ = strconv.Atoi(v)
			uiDs = append(uiDs, uid)
		}
		user, err := dao.GetUserInfoUidKey(uiDs)
		if err != nil {
			return likeList, err
		}
		for _, v := range val {
			uid, _ = strconv.Atoi(v)
			likeList["list"] = append(likeList["list"], map[string]interface{}{
				"uid":        uid,
				"p_id":       p.Pid,
				"first_name": user[uid].(models.User).FirstName,
				"last_name":  user[uid].(models.User).LastName,
				"avatar":     user[uid].(models.User).Avatar,
			})
		}
		return likeList, nil
	}
	return likeList, nil
}
