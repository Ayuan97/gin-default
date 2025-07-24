package main

import (
	"fmt"
	"justus/internal/models"
	"justus/pkg/logger"
	"justus/pkg/setting"
	"justus/routers"
	"log"
)

func init() {
	setting.Setup()
	logger.Setup()
	models.Setup()
}

func main() {
	log.Printf("🚀 启动 Justus API 服务，端口: %d", setting.ServerSetting.HttpPort)
	router := routers.InitRouter()
	router.Run(fmt.Sprintf(":%d", setting.ServerSetting.HttpPort))
}
