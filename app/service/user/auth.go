package user

import (
	"fmt"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	//db := c.Value("DB").(*gorm.DB)

	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Println(username, password)

	r := result.New()
	//校验密码
	jwtConfig := config.Cfg.Jwt
	jwtToken, err := authUtil.GenerateJwtToken(jwtConfig, 1)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail("登录失败"))
		return
	}
	fmt.Println(jwtToken)

	//生成散列hash 并存到redis
	token, err := authUtil.SaveMd5Token(jwtToken)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail("登录失败"))
		return
	}
	fmt.Println(token)
	c.SetCookie("token", token, jwtConfig.Expires, "/", c.Request.Host, false, true)
	r.SetData(gin.H{"token": token})

	c.JSON(http.StatusOK, r)
}

func Logout(c *gin.Context) {

}
