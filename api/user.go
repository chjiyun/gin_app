package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Users(c *gin.Context) {
	// 你也可以使用一个结构体
	var msg struct {
		Name    string
		Message string
		Age     int
	}
	msg.Name = "Chjiyun"
	msg.Message = "hey, Hello World!"
	msg.Age = 27

	logrus.Info(msg)
	c.JSON(http.StatusOK, msg)
}
