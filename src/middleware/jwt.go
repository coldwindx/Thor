package middleware

import (
	"Thor/common"
	"Thor/config"
	"Thor/ctx"
	"Thor/src/services"
	"Thor/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func JWTAuth(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if "" == auth {
			common.Fail(c, common.Errors.TokenError.Code, "请先登录")
			c.Abort()
			return
		}
		auth = auth[len(services.TokenType)+1:]
		// token解析校验
		token, err := jwt.ParseWithClaims(auth, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.Jwt.Secret), nil
		})
		if nil != err {
			common.Fail(c, common.Errors.TokenError.Code, common.Errors.TokenError.Message)
			c.Abort()
			return
		}
		// 黑名单校验
		if services.JwtService.IsInBlackList(auth) {
			common.Fail(c, common.Errors.TokenError.Code, common.Errors.TokenError.Message)
			c.Abort()
			return
		}
		claims := token.Claims.(*services.CustomClaims)
		// token发布者校验
		if GuardName != claims.Issuer {
			common.Fail(c, common.Errors.TokenError.Code, common.Errors.TokenError.Message)
			c.Abort()
			return
		}
		// token续签
		if claims.ExpiresAt-time.Now().Unix() < config.Config.Jwt.RefreshGracePeriod {
			lock := tools.Lock("refresh_token_lock", config.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				err, user := services.JwtService.GetUserInfo(GuardName, claims.Id)
				if nil != err {
					ctx.Logger.Error(err.Error())
					lock.Release()
				} else {
					newToken, _, _ := services.JwtService.CreateToken(GuardName, user)
					c.Header("new-token", newToken.AccessToken)
					c.Header("new-expires-in", strconv.Itoa(newToken.ExpiresIn))
					_ = services.JwtService.JoinBlackList(token)
				}
			}
		}
		c.Set("token", token)
		c.Set("id", claims.Id)
	}
}
