package api_require

import (
	"github.com/gin-gonic/gin"
)


func Common() gin.HandlerFunc {
	return func(c *gin.Context) {

		version := c.GetHeader("version")       //版本号
		uuid := c.GetHeader("uuid")             //设备号
		deviceType := c.GetHeader("deviceType") //设备类型
		lange := c.GetHeader("lange")           //语言
		timeZone := c.GetHeader("timeZone")     //时区
		c.Set("version",version)
		c.Set("uuid",uuid)
		c.Set("device_type", deviceType)
		c.Set("lange",lange)
		c.Set("time_zone", timeZone)
		c.Next()
	}
}
