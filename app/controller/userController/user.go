package userController

import (
	"gin_app/app/controller/userController/userVo"
	"gin_app/app/result"
	"gin_app/app/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	r := userService.Login(c)
	c.JSON(http.StatusOK, r)
}

func Register(c *gin.Context) {
	r := userService.Register(c)
	c.JSON(http.StatusOK, r)
}

func Logout(c *gin.Context) {
	r := userService.Logout(c)
	c.JSON(http.StatusOK, r)
}

func ResetPassword(c *gin.Context) {
	r := userService.ResetPassword(c)
	c.JSON(http.StatusOK, r)
}

func GetPageUsers(c *gin.Context) {
	var reqVo userVo.UserPageReqVo
	r := c.Value("Result").(*result.Result)
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	userService.GetPageUsers(c, reqVo)
	c.JSON(http.StatusOK, r)
}

func GetCurrentUser(c *gin.Context) {
	r := c.Value("Result").(*result.Result)
	userService.GetCurrentUser(c)
	c.JSON(http.StatusOK, r)
}
