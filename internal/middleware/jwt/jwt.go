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
		var code int

		code = e.SUCCESS
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
				userId, _ := strconv.Atoi(claims.Subject)
				c.Set("userId", userId)
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
