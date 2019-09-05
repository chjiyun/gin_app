package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	c.JSON(http.StatusOK, msg)
}
