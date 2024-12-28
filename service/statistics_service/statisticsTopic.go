package statistics_service

import (
	"encoding/json"
	"fmt"
	"justus/global"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"strconv"
	"strings"
)

//处理数据
func StatisticsTopicHandle()bool  {
	var data Data
	var idsList []string
	var topicMap  = make(map[int]map[string]int)
	isRun := true
	for{
		if isRun == false{
			break
		}
		key := rediskey.StatisticsTopicList
		value, err := gredis.RPop(key)
		if value != ""{		//取数据 并把数据入map中
			err = json.Unmarshal([]byte(value), &data)
			if err != nil {
				return false
			}
			idsList = strings.Split(data.Ids,",")
			for _,v := range idsList{
				id, _ := strconv.Atoi(v)
				//自增计数
				if _,ok := topicMap[id];ok{
					if _,ok = topicMap[id][data.DataType];ok{
						if data.Inc == 1{
							topicMap[id][data.DataType]++
						}
						if data.Inc == -1{
							topicMap[id][data.DataType]--
						}
					}else{
						if data.Inc == 1{
							topicMap[id][data.DataType] = 1
						}
						if data.Inc == -1{
							topicMap[id][data.DataType] = -1
						}

					}
				}else{
					if data.Inc == 1{
						topicMap[id] = map[string]int{data.DataType:1}
					}
					if data.Inc == -1{
						topicMap[id] = map[string]int{data.DataType:-1}
					}
				}
			}
		}else{		//队列取结束 更新到库

			var weightIdList  []int  //需要计算权重的
			var initDataIdList []int //需要初始化数据的
			if len(topicMap) > 0 {
				fmt.Println("topicMap",topicMap)
				topicStatisticsModel := models.TopicStatistics{}
				for id, v := range topicMap {
					//更新所有计数
					topicStatisticsModel.TopicId = id
					topicStatisticsModel.IncData(v)
					curTmpData := getCurColumnData(id, FromTopic)
					for v_key, num := range v {
						if v_key == Impression { //曝光数
							curTmpImpression := curTmpData.Impression + num
							if curTmpImpression >= startImpression && curTmpImpression <= endImpression { //需要计算权重
								weightIdList = append(weightIdList, id)
							} else if curTmpImpression > endImpression { //初始化权重
								initDataIdList = append(initDataIdList, id)
							}
							curTmpData.Impression = curTmpImpression
						}
						if v_key == Like {
							curTmpData.Like = curTmpData.Like + num
						}
						if v_key == Comment {
							curTmpData.Comment = curTmpData.Comment + num
						}
						if v_key == FollowPhoto {
							curTmpData.FollowPhoto = curTmpData.FollowPhoto + num
						}
						_ = setCurColumnData(id, FromTopic, *curTmpData)
					}
				}


				//更新权重
				if len(weightIdList) > 0 {
					fmt.Println("weightlist", weightIdList)
					topicModel := models.Topics{}
					for _, v := range weightIdList {
						topicModel.Weight = getWeight(v, FromTopic)
						topicModel.ID = v
						fmt.Println("topicmodel", topicModel)
						err = topicModel.UpdateTopicWeight()
						if err != nil {
							global.Logger.Error("UpdateTopicWeight failed", v)
						}
					}
				}
				//需要初始化权重的
				if len(initDataIdList) > 0 {
					fmt.Println("initDataIdList", initDataIdList)
					topicStatisticsModel2 := models.TopicStatistics{}
					err = topicStatisticsModel2.InitWeightData(initDataIdList)
					if err != nil {
						global.Logger.Error("initWeightData failed", initDataIdList)
					}
					for _, v := range initDataIdList {
						_ = setCurColumnData(v, FromTopic, curColumnData{1, 1, 1, 1})
					}
				}
			}
			isRun = false
		}

	}
	return true

}

