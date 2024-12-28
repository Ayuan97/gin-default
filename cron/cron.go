package main

import (
	"gin-default/models"
	"gin-default/pkg/gredis"
	"gin-default/pkg/logger"
	"gin-default/pkg/setting"
	"gin-default/pkg/util"
	"time"
)

func init() {
	setting.Setup()
	logger.Setup()
	models.Setup()
	util.Setup()
	gredis.Setup()
}

func main() {
	var err error
	c := cron.New(cron.WithSeconds())
	////添加2秒钟定时任务 处理
	//_, err = c.AddFunc("*/2 * * * * *", func() {
	//	fmt.Println("start_time_topic:", time.Now())
	//	statistics_service.StatisticsTopicHandle()
	//	fmt.Println("end_time_topic:", time.Now())
	//	fmt.Println("start_time_picture:", time.Now())
	//	statistics_service.StatisticsPictureHandle()
	//	fmt.Println("end_time_picture:", time.Now())
	//})
	//if err != nil {
	//	fmt.Println("AddFun:", err)
	//	return
	//}
	//
	////每天删除过期的推荐信息
	//_, err = c.AddFunc("0 0 0 * * *", func() {
	//	fmt.Println("开始删除过期推荐信息:", time.Now())
	//	recommend_service.DeleteRecommendInfo()
	//	fmt.Println("删除过期推荐信息结束:", time.Now())
	//})
	//if err != nil {
	//	fmt.Println("AddFun:DeleteRecommendInfo:", err)
	//	return
	//}

	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)

		}
	}
}
