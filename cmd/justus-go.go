package main

import (
	"fmt"
	"log"

	"justus/internal/models"
	"justus/internal/routers"
	"justus/pkg/gredis"
	"justus/pkg/logger"
	"justus/pkg/setting"
)

func init() {
	setting.Setup()
	logger.Setup()
	gredis.Setup()
	models.Setup()
}

func main() {
	log.Printf("🚀 启动 Justus API 服务，端口: %d", setting.ServerSetting.HttpPort)

	// 使用依赖注入初始化路由
	router, err := routers.InitRouterWith()
	if err != nil {
		log.Fatalf("❌ 初始化路由失败: %v", err)
	}

	log.Println("依赖注入系统初始化完成")
	router.Run(fmt.Sprintf(":%d", setting.ServerSetting.HttpPort))
}
