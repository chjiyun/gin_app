package service

import (
	"gin_app/app/common"
	"gin_app/app/model"
	"gin_app/app/util"
	"gin_app/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Upload 接收上传的文件
func Upload(c *gin.Context) {
	result := &common.Result{}
	db := c.Value("DB").(*gorm.DB)

	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, result.Fail("上传失败", err.Error()))
		return
	}

	ext := filepath.Ext(f.Filename)
	mimetype := f.Header["Content-Type"][0]
	var filetype string
	if len(mimetype) > 0 {
		i := strings.Index(mimetype, "/")
		if i > 0 {
			filetype = mimetype[:i]
		}
	}
	uid := uuid.NewV4().String()
	localFilename := uid + ext
	year, month, _ := time.Now().Date()
	relativePath := filepath.Join("files", strconv.Itoa(year), strconv.Itoa(int(month)), localFilename)
	sourcepath := filepath.Join(config.Cfg.Basedir, relativePath)
	// 文件所在文件夹目录
	dirname := filepath.Dir(sourcepath)
	err = os.MkdirAll(dirname, 0666)
	if err != nil {
		c.JSON(200, result.Fail("服务器内部错误", err.Error()))
		return
	}

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(f, sourcepath)
	if err != nil {
		c.JSON(200, result.Fail("上传失败", err.Error()))
		return
	}

	file := model.File{
		Name:     f.Filename,
		Uid:      uid,
		Ext:      ext[1:],
		Type:     filetype,
		MimeType: mimetype,
		Path:     relativePath,
		Size:     uint(f.Size),
	}
	res := db.Create(&file)
	if res.Error != nil {
		c.JSON(200, result.Fail("上传失败", res.Error.Error()))
		return
	}
	c.JSON(200, result.Success("", gin.H{
		"uid": file.Uid, "ext": file.Ext, "name": file.Name, "size": file.Size,
	}))
}

// Download 下载文件
func Download(c *gin.Context) {
	result := &common.Result{}
	id := c.Param("id")
	db := c.Value("DB").(*gorm.DB)
	var file model.File

	res := db.First(&file, "uid", id)
	if res.Error != nil {
		c.JSON(200, result.Fail("文件不存在", res.Error.Error()))
		return
	}
	sourcePath := filepath.Join(config.Cfg.Basedir, file.Path)
	if !util.CheckFileIsExist(sourcePath) {
		c.JSON(200, result.Fail("文件不存在", nil))
		return
	}
	c.File(sourcePath)
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
