package statistics_service

import (
	"encoding/json"
	"fmt"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"strconv"
)
const startImpression = 100	//开始计算权重的曝光数
const endImpression = 200	//结束计算权重的曝光数
const baseWeight = 0.25  //初始权重

//当前实时数据的结构体
type curColumnData struct {
	Impression  int `json:"impression"`
	Like        int `json:"like"`
	Comment     int `json:"comment"`
	FollowPhoto int `json:"follow_photo"`
}



//获取当前的实时数据缓存
func getCurColumnData(id int,fromType int)*curColumnData{
	curData := curColumnData{0,0,0,0}
	var key string
	if fromType == FromTopic {
		key = rediskey.GetStatisticsTopicCurDataKey(id)
	}else if fromType == FromPicture {
		key = rediskey.GetStatisticsPictureCurDataKey(id)
	}else{
		return nil
	}
	value := gredis.Get(key)
	if value!=""{
		err := json.Unmarshal([]byte(value), &curData)
		if err != nil {
			return nil
		}
	}
	return &curData
}

//设置当前实时数据缓存
func setCurColumnData(id int,fromType int,data curColumnData)error{
	var key string
	if fromType == FromTopic {
		key = rediskey.GetStatisticsTopicCurDataKey(id)
	}else if fromType == FromPicture {
		key = rediskey.GetStatisticsPictureCurDataKey(id)
	}else{
		var err error
		return err
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = gredis.Set(key,marshal,0)
	if err != nil {
		return err
	}
	return nil
}

//获取权重
func getWeight(id int,fromType int)float64{
	curData := getCurColumnData(id,fromType)
	if curData.Impression > 0{
		weight := float64(curData.FollowPhoto)/float64(curData.Impression) + 0.2*float64(curData.Like)/float64(curData.Impression) + 0.3*float64(curData.Comment)/float64(curData.Impression)
		floatStr := fmt.Sprintf("%.6f", weight)
		inst, _ := strconv.ParseFloat(floatStr, 64)
		//fmt.Println("weight",weight)
		return inst
	}
	return 0
}



