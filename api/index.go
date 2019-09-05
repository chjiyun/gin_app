package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	// 你也可以使用一个结构体
	var msg struct {
		Name    string `json:"user"`
		Message string
		Number  int
	}
	msg.Name = "Lena"
	msg.Message = "hey, Hello World!"
	msg.Number = http.StatusOK
	c.JSON(http.StatusOK, msg)
}
