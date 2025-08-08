package jwt

import (
	"errors"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/pkg/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int = e.SUCCESS
		token := c.GetHeader("Authorization")
		token = strings.Replace(token, "Bearer ", "", -1)
		token = strings.Replace(token, "Bearer", "", -1)
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
			if claims != nil {
				// 兼容旧字段 Subject，同时支持新字段 AdminUserID
				userId := 0
				if claims.AdminUserID != 0 {
					userId = claims.AdminUserID
				} else {
					// 兼容旧：Subject 存储用户ID
					if claims.Subject != "" {
						userId, _ = strconv.Atoi(claims.Subject)
					}
				}
				if userId != 0 {
					c.Set("userId", userId)
				}
				// 多租户注入
				if claims.IsSuper {
					c.Set("isSuper", true)
				}
				if claims.TenantID != 0 {
					c.Set("tenantId", claims.TenantID)
				}
				if len(claims.TenantIDs) > 0 {
					c.Set("tenantIds", claims.TenantIDs)
				}
			}

		}

		if code != e.SUCCESS {
			appG := app.Gin{C: c}
			appG.Unauthorized(code)
			c.Abort()
			return
		}

		c.Next()
	}
}
