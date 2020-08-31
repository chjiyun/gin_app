package test

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// String 数组和切片
func String(c *gin.Context) {
	// 一个 UTF8 编码的字符可能会占多个字节，比如汉字就需要 3~4 个字节来存储
	x := "hello"
	// 将 string 转为 rune slice（此时 1 个 rune 可能占多个 byte）
	xRunes := []rune(x)
	xRunes[0] = '我'
	x = string(xRunes)

	c.JSON(http.StatusOK, x)
}
