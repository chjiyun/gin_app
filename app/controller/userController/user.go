package userController

import (
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

func GetUsers(c *gin.Context) {
	r := userService.GetUsers(c)
	c.JSON(http.StatusOK, r)
}
