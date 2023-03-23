package userController

import (
	"gin_app/app/controller/userController/userIpVo"
	"gin_app/app/controller/userController/userVo"
	"gin_app/app/result"
	"gin_app/app/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	r := result.New()
	var reqVo userVo.UserLoginReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := userService.Login(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func Register(c *gin.Context) {
	r := result.New()
	var reqVo userVo.UserRegisterReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	err := userService.Register(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r)
}

func Logout(c *gin.Context) {
	r := result.New()
	userService.Logout(c)
	c.JSON(http.StatusOK, r)
}

func ResetPassword(c *gin.Context) {
	r := result.New()
	var reqVo userVo.UserResetPasswordReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	flag, err := userService.ResetPassword(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(flag))
}

func GetUserPage(c *gin.Context) {
	r := result.New()
	var reqVo userVo.UserPageReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, _ := userService.GetUserPage(c, reqVo)
	c.JSON(http.StatusOK, r.Success(res))
}

func GetCurrentUser(c *gin.Context) {
	r := result.New()
	res := userService.GetCurrentUser(c)
	c.JSON(http.StatusOK, r.Success(res))
}

func GetUserIpPage(c *gin.Context) {
	r := result.New()
	var reqVo userIpVo.UserIpPageReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, _ := userService.GetUserIpPage(c, reqVo)
	c.JSON(http.StatusOK, r.Success(res))
}
