package api

import (
	"fmt"
	"gin_app/util"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
)

// Upload 接收上传的文件
func Upload(c *gin.Context) {
	// formType := c.PostForm("type")
	// single file
	file, _ := c.FormFile("file")

	filename := filepath.Join("files", file.Filename)
	filetype := filepath.Ext(file.Filename)
	filetype = string([]rune(filetype)[1:])
	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, filename)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	path, _ := os.Getwd()
	path = filepath.Join(path, filename)

	c.JSON(http.StatusOK, util.ResponseMsg{Code: 200, Msg: "success", Data: gin.H{"filepath": path, "filetype": filetype}})
}

// Download 下载文件
func Download(c *gin.Context) {
	id := c.Param("id")
	path := filepath.Join("files", id)
	fmt.Println(path)
	if !CheckFileIsExist(path) {
		c.String(http.StatusOK, "file not found")
		return
	}
	c.File(path)
}

// ExtractWord 提取 word文档
func ExtractWord(c *gin.Context) {
	file, _ := c.FormFile("file")
	f, _ := file.Open()
	defer f.Close()
	r, err := docx.ReadDocxFromMemory(f, file.Size)
	if err != nil {
		panic(err)
	}
	docx1 := r.Editable()
	text := docx1.GetContent()
	// fmt.Println(text)
	r.Close()
	c.JSON(200, text)
}
