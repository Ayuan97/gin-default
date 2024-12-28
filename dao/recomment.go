package dao

import (
	"github.com/go-redis/redis/v8"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"time"
)

func GetRecommendList(NotPid []string, offset int, limit int, lange string) ([]*models.PictureRecommend, error) {

	pictureRecommend := models.PictureRecommend{}
	return pictureRecommend.GetPictureRecommendListNot(NotPid, offset, limit, lange)

}

// GetRecommendListReadRecord 获取推荐列表已读记录
func GetRecommendListReadRecord(userId int) []string {
	//获取已读记录
	key := rediskey.GetRecommendListReadRecordKey(userId)
	val, _ := gredis.Zrange(key, 0, -1)
	return val
}

// AddRecommendListReadRecord 添加推荐列表已读记录
func AddRecommendListReadRecord(userId int, pid []string) {
	//添加已读记录
	key := rediskey.GetRecommendListReadRecordKey(userId)
	//当前时间
	now := time.Now().Unix()
	for _, v := range pid {
		gredis.Zadd(key, redis.Z{Score: float64(now), Member: v})
	}
}

// GetLanguageRecommendCount 获取每个语言的推荐数量
func GetLanguageRecommendCount() ([]*models.LangeRecommendCount, error) {
	pictureRecommend := models.PictureRecommend{}
	return pictureRecommend.GetPictureRecommendCount()
}

// 查询20001个信息
func GetRecommendDeleteInfo(lange string, offset int) bool {
	pictureRecommend := models.PictureRecommend{
		Lange: lange,
	}
	return pictureRecommend.GetPictureRecommendExpire(offset)
}
