package controller

import (
	"gin_app/app/result"
	userService "gin_app/app/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	r := result.New()
	userService.Login(c, r)
	c.JSON(http.StatusOK, r)
}

func Register(c *gin.Context) {
	r := result.New()
	userService.Register(c, r)
	c.JSON(http.StatusOK, r)
}

func Logout(c *gin.Context) {
	r := result.New()
	userService.Logout(c)
	c.JSON(http.StatusOK, r)
}

func ResetPassword(c *gin.Context) {
	r := result.New()
	userService.ResetPassword(c, r)
	c.JSON(http.StatusOK, r)
}
