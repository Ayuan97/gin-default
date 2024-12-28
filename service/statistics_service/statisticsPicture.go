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
func StatisticsPictureHandle()bool  {
	var data Data
	var idsList []string
	var pictureMap = make(map[int]map[string]int)
	var uidMap = make(map[int]int)
	isRun := true
	for{
		if isRun == false{
			break
		}
		key := rediskey.StatisticsPictureList
		value, err := gredis.RPop(key)
		if value != ""{		//取数据 并把数据入map中
			err = json.Unmarshal([]byte(value), &data)
			if err != nil {
				return false
			}
			idsList = strings.Split(data.Ids,",")
			uidsList := strings.Split(data.Uids,",")
			uidMap = handleUidsMap(uidsList,uidMap)				//用户曝光数统计
			pictureMap = handlePictureIdsMap(idsList,data,pictureMap)	//照片 各个类型数据统计
		}else{		//队列取结束 更新到库
			if len(pictureMap) > 0 {
				fmt.Println("pictureMap", pictureMap)
				var weightIdList []int   //需要计算权重的
				var initDataIdList []int //需要初始化数据的
				pictureStatisticsModel := models.PictureStatistics{}
				for id, v := range pictureMap {
					//更新所有计数
					pictureStatisticsModel.PId = id
					pictureStatisticsModel.IncData(v)
					curTmpData := getCurColumnData(id, FromPicture)
					for vKey, num := range v {
						if vKey == Impression { //曝光数
							curTmpImpression := curTmpData.Impression + num
							if curTmpImpression >= startImpression && curTmpImpression <= endImpression { //需要计算权重
								weightIdList = append(weightIdList, id)
							} else if curTmpImpression > endImpression { //初始化权重
								initDataIdList = append(initDataIdList, id)
							}
							curTmpData.Impression = curTmpImpression
						}
						if vKey == Like {
							curTmpData.Like = curTmpData.Like + num
						}
						if vKey == Comment {
							curTmpData.Comment = curTmpData.Comment + num
						}
						if vKey == FollowPhoto {
							curTmpData.FollowPhoto = curTmpData.FollowPhoto + num
						}
						_ = setCurColumnData(id, FromPicture, *curTmpData)
					}
				}


				//更新权重
				if len(weightIdList) > 0 {
					fmt.Println("picture_weightlist", weightIdList)
					pictureRecommendModel := models.PictureRecommend{}
					topicPictureModel := models.TopicPicture{}
					for _, v := range weightIdList {
						pictureRecommendModel.Weight = getWeight(v, FromPicture)
						pictureRecommendModel.PId = v
						err = pictureRecommendModel.UpdatePictureWeight()
						if err != nil {
							global.Logger.Error("UpdatePictureWeight_recommend failed", v)
						}

						topicPictureModel.Weight = getWeight(v, FromPicture)
						topicPictureModel.PId = v
						err = topicPictureModel.UpdatePictureWeight()
						if err != nil {
							global.Logger.Error("UpdatePictureWeight_rtopic failed", v)
						}
					}
				}
				//需要初始化权重的
				if len(initDataIdList) > 0 {
					fmt.Println("picture_initDataIdList", initDataIdList)
					PictureStatisticsModel2 := models.PictureStatistics{}
					err = PictureStatisticsModel2.InitWeightData(initDataIdList)
					if err != nil {
						global.Logger.Error("initWeightData failed", initDataIdList)
					}
					for _, v := range initDataIdList {
						_ = setCurColumnData(v, FromPicture, curColumnData{1, 1, 1, 1})
					}
				}
			}

			if len(uidMap) > 0{
				for id, num := range uidMap {
					userPageModel := models.UserPage{}
					userPageModel.Uid = id
					_ = userPageModel.IncUserHotNum(num)
				}
			}
			isRun = false
		}

	}
	return true

}


func handleUidsMap(uid []string,uidsMap map[int]int) map[int]int  {
	for _,v := range uid{
		id, _ := strconv.Atoi(v)
		//自增计数
		if _,ok := uidsMap[id];ok{
			if _,ok = uidsMap[id];ok{
				uidsMap[id]++
			}else{
				uidsMap[id] = 1
			}
		}else{
			uidsMap[id] = 1
		}
	}
	return uidsMap

}

func handlePictureIdsMap(idsList []string,data Data,pictureMap map[int]map[string]int) map[int]map[string]int  {
	for _,v := range idsList{
		id, _ := strconv.Atoi(v)
		//自增计数
		if _,ok := pictureMap[id];ok{
			if _,ok = pictureMap[id][data.DataType];ok{
				if data.Inc == 1 {
					pictureMap[id][data.DataType]++
				}
				if data.Inc == -1 {
					pictureMap[id][data.DataType]--
				}
			}else{
				if data.Inc == 1 {
					pictureMap[id][data.DataType] = 1
				}
				if data.Inc == -1 {
					pictureMap[id][data.DataType] = -1
				}
			}
		}else{
			if data.Inc == 1 {
				pictureMap[id] = map[string]int{data.DataType: 1}
			}
			if data.Inc == -1 {
				pictureMap[id] = map[string]int{data.DataType: -1}
			}

		}
	}
	return pictureMap

}
