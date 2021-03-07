package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload 接收上传的文件
func Upload(c *gin.Context) {
	x := "test upload..."
	c.JSON(http.StatusOK, x)
}

// Download 下载文件
func Download(c *gin.Context) {
	x := "test download...1"
	c.JSON(http.StatusOK, x)
}

// ExtractWord 提取 word文档
func ExtractWord(c *gin.Context) {
	c.JSON(http.StatusOK, "哈哈哈哈哈哈哈12")
}
