package tenant

import (
	"github.com/gin-gonic/gin"
)

// Resolve 解析并注入租户上下文
// 方案：严格使用 JWT 中的 tenantId，不再允许通过 Header/子域名切换
func Resolve() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果 JWT 已注入 tenantId，直接放行；否则交给后续中间件/处理函数报错
		c.Next()
	}
}
