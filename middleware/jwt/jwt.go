package jwt

import (
	"justus/pkg/e"
	"justus/pkg/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//type UserJwt struct {
//	UserId int
//	ExpiresAt int
//}
//var User = &UserJwt{}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("Authorization")
		token = strings.Replace(token, "Bearer ", "", -1)
		token = strings.Replace(token, "Bearer", "", -1)
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
			if claims != nil {
				userId,_ := strconv.Atoi(claims.Subject)
				//User.ExpiresAt = int(claims.ExpiresAt)
				c.Set("userId",userId)
				//fmt.Println("token:",token,"claims",claims)
			}

		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
