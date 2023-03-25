package middleware

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := result.New()
		jwtConfig := config.Cfg.Jwt

		reWhitelist, _ := regexp.Compile(jwtConfig.Whitelist)
		if reWhitelist.MatchString(c.Request.URL.Path) {
			return
		}

		token, err := c.Cookie("token")
		if err != nil && errors.Is(err, http.ErrNoCookie) {
			token = c.Query("token")
		}
		// 未登录或cookie已过期
		if token == "" {
			c.JSON(http.StatusOK, r.FailType(common.UnLogin).SetCode(401))
			c.Abort()
			return
		}
		//从redis取真正的token
		jwtToken, err := config.RedisDb.Get(c, token).Result()
		if err != nil {
			c.JSON(http.StatusOK, r.FailType(common.UnLogin).SetCode(401))
			c.Abort()
			return
		}

		claims, err := authUtil.ParseJwtToken(jwtToken, jwtConfig.SecretKey)
		if err != nil {
			r.FailType(common.UnLogin).SetCode(401)

			// 过期就续签token
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = authUtil.RenewJwtToken(jwtConfig, claims.UserId, token)
				if err != nil {
					c.JSON(http.StatusOK, r)
					c.Abort()
					return
				}
				splitHost := strings.Split(c.Request.Host, ":")
				if len(splitHost) > 0 {
					c.SetCookie("token", token, jwtConfig.Expires, "/", splitHost[0], false, true)
				}
			} else {
				c.JSON(http.StatusOK, r)
				c.Abort()
				return
			}
		}

		//可以在这里做权限校验

		//将解析出的信息传递下去
		c.Set("session", claims)
		c.Next()
	}
}
