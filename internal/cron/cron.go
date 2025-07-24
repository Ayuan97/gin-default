package main

import (
	"fmt"
	"justus/internal/models"
	"justus/pkg/gredis"
	"justus/pkg/logger"
	"justus/pkg/setting"
	"justus/pkg/util"
	"time"

	"github.com/robfig/cron/v3"
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
	//添加2秒钟定时任务 处理
	_, err = c.AddFunc("*/2 * * * * *", func() {
		fmt.Println("start_time_topic:", time.Now())
	})
	if err != nil {
		fmt.Println("AddFun:", err)
		return
	}

	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)

		}
	}
}
