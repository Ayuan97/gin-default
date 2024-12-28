package follow_service

import (
	"github.com/go-redis/redis/v8"
	"justus/dao"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/pkg/setting"
	"justus/service"
	"justus/service/picture_library_service"
	"strconv"
	"strings"
	"time"
)

type PostFollowList struct {
	Page int `json:"page"`
	UId  int `json:"uid"`
}

// GetFollowList 获取关注列表
func GetFollowList(uid int, page int, userId int) *service.ImgList {
	var data service.ImgList
	//是否有关注列表集合
	isFollowListMap := IsFollowListMap(uid)
	var MapNumber int64
	userFollowMap := GetUserFollowMap(uid)
	userCollectTopicMap := GetUserCollectTopicMap(uid)
	mapList := append(userFollowMap, userCollectTopicMap...)

	if isFollowListMap > 0 {
		MapNumber = MergeMap(mapList, uid)
		AddUserFollowMapUpdateTime(uid)
	} else {
		//没有关注列表集合
		//合并有序集合
		MapNumber = MergeMap(mapList, uid)
		AddUserFollowMapUpdateTime(uid)
	}
	if MapNumber > 0 {
		var offset int
		if page > 0 {
			offset = (page - 1) * setting.AppSetting.PageSize
		} else {
			offset = 0
		}
		//有序集合分页
		key := rediskey.GetFollowListMapKey(uid)
		value, _ := gredis.Zrevrange(key, int64(offset), int64(offset+setting.AppSetting.PageSize-1))
		var pid []int
		if len(value) > 0 {
			for _, v := range value {
				id, err := strconv.Atoi(v)
				if err == nil {
					pid = append(pid, id)
				}
			}
			imgListData, _ := dao.GetPictureLibrarys(pid)
			imgList, err := picture_library_service.PictureLibraryList(imgListData, userId, 0)
			if err != nil {
				data.ImgList = imgList
				return &data
			}
			data.ImgList = imgList
			return &data
		}
	}
	list := make([]*models.PictureLibraryList, 0)
	data.ImgList = list
	return &data
}

func MergeMap(mapList []string, uid int) int64 {
	//合并有序集合
	key := rediskey.GetFollowListMapKey(uid)
	var z *redis.ZStore
	z = new(redis.ZStore)
	z.Aggregate = "Max"
	for _, v := range mapList {
		//查看 集合名称前缀是否包含justus
		if strings.Contains(v, setting.RedisSetting.Prefix) {
			z.Keys = append(z.Keys, v)
		} else {
			z.Keys = append(z.Keys, setting.RedisSetting.Prefix+v)
		}
	}
	if len(z.Keys) > 0 {
		value, _ := gredis.ZunionStore(key, z)
		if value > 2000 {
			delCount := value - 2000
			_, _ = gredis.Zremrangebyrank(key, 0, delCount)
		}
		return value
	}
	return 0
}

//IsFollowListMap 获取关注列表map
func IsFollowListMap(uid int) int64 {
	key := rediskey.GetFollowListMapKey(uid)
	value, _ := gredis.Zcard(key)
	return value
}

// GetUserFollowMap 获取用户关注的用户map
func GetUserFollowMap(uid int) []string {
	key := rediskey.GetUserFollowUserKey(uid)
	value := gredis.SMembers(key)
	var keysMap []string
	if len(value) > 0 {
		for _, v := range value {
			id, _ := strconv.Atoi(v)
			keysMap = append(keysMap, setting.RedisSetting.Prefix+rediskey.GetUserPictureMapKey(id))
		}

	}
	return keysMap

}

// GetUserCollectTopicMap 获取用户收藏的话题map
func GetUserCollectTopicMap(uid int) []string {
	key := rediskey.GetUserCollectTopicKey(uid)
	value := gredis.SMembers(key)
	var keysMap []string
	if len(value) > 0 {
		for _, v := range value {
			id, _ := strconv.Atoi(v)
			keysMap = append(keysMap, setting.RedisSetting.Prefix+rediskey.GetTopicPictureMapKey(id))
		}

	}
	return keysMap
}

// AddUserFollowMapUpdateTime 关注列表map的更新时间
func AddUserFollowMapUpdateTime(uid int) bool {
	key := rediskey.GetFollowListMapUpdateTime(uid)
	err := gredis.Set(key, float64(time.Now().Unix()), time.Hour*24*60)
	if err != nil {
		return false
	}
	return true
}
