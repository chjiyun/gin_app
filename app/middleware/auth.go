package middleware

import (
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
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
		if err != nil {
			r.SetResult(result.ResultMap["illegalVisit"], "token缺失")
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}
		//从redis取真正的token
		realToken, err := config.RedisDb.Get(c, token).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, r.Fail("token过期"))
			c.Abort()
			return
		}

		claims, err := authUtil.ParseJwtToken(realToken, jwtConfig.SecretKey)
		if err != nil {
			r.SetResult(result.ResultMap["illegalVisit"], "")
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.RegisteredClaims.ExpiresAt.Unix() {
			r.SetResult(result.ResultMap["illegalVisit"], "token过期")
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}

		//可以在这里做权限校验

		//将解析出的信息传递下去
		c.Set("user", claims)
		c.Next()
	}
}
