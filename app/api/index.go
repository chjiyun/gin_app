package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Index 测试语法的 API
func Index(c *gin.Context) {
	// 你也可以使用一个结构体
	var res struct {
		Msg  string
		Code int
		Date string
	}
	x := 1234567890
	res.Msg = fmt.Sprintf("x=%d", x)
	res.Code = 200
	res.Date = time.Now().Format("2006-01-02")
	c.JSON(http.StatusOK, res)
}
