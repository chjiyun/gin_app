package middleware

import (
	"errors"
	"fmt"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"net/http"
	"regexp"

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
			r.SetResult(result.ResultMap["unLogin"], "")
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}
		//从redis取真正的token
		jwtToken, err := config.RedisDb.Get(c, token).Result()
		if err != nil {
			r.SetResult(result.ResultMap["unLogin"], "")
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}

		claims, err := authUtil.ParseJwtToken(jwtToken, jwtConfig.SecretKey)
		if err != nil {
			fmt.Println(err)
			r.SetResult(result.ResultMap["unLogin"], "")

			// 过期就续签token
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = authUtil.RenewJwtToken(jwtConfig, claims.UserId, token)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusUnauthorized, r)
					c.Abort()
					return
				}
			}
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}

		//可以在这里做权限校验

		//将解析出的信息传递下去
		c.Set("session", claims)
		c.Next()
	}
}
