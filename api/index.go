package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index 测试语法的 API
func Index(c *gin.Context) {
	// 你也可以使用一个结构体
	var msg struct {
		Result string
		Number string
	}
	x := 1234567890
	msg.Result = "test Float"
	msg.Number = fmt.Sprintf("x=%d", x)
	c.JSON(http.StatusOK, msg)
}
