package recovers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"justus/internal/global"
	"runtime"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\033[1;31;40m%s\033[0m\n", "系统错误:"+fmt.Sprint(err))
				errInfo := logrus.Fields{}
				for i := 2; i <= 4; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					//格式化输出 红色警告
					fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("%s:%d", file, line))
					errInfo[fmt.Sprintf("%d", i)] = fmt.Sprintf("%s:%d", file, line)

				}
				//errInfo["error"] = string(debug.Stack()) //记录全部信息
				global.Logger.WithFields(errInfo).Error("错误:", err, "\n", "错误位置:", errInfo["2"])
				//global.Logger.Error("捕获异常:", err)
			}

		}()
		c.Next()
	}
}
