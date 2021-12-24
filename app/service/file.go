package service

import (
	"gin_app/app/common"
	"gin_app/app/model"
	"gin_app/app/util"
	"gin_app/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
	uuid "github.com/satori/go.uuid"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

// Upload 接收上传的文件
func Upload(c *gin.Context) {
	result := &common.Result{}
	db := c.Value("DB").(*gorm.DB)

	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, result.Fail("上传失败", err))
		return
	}

	ext := filepath.Ext(f.Filename)
	// mimetype := f.Header["Content-Type"][0]
	mfile, _ := f.Open()
	mime, err := mimetype.DetectReader(mfile)
	if err != nil {
		c.JSON(200, result.Fail("MIME type detect failed", err))
		return
	}
	mtype := mime.String()
	var filetype string

	if i := strings.Index(mtype, "/"); i > 0 {
		filetype = mtype[:i]
	}
	uid := uuid.NewV4().String()
	snowId := idgen.NextId()
	localName := util.ToString(snowId) + ext
	year, month, _ := time.Now().Date()
	relativePath := filepath.Join("files", util.ToString(year), util.ToString(int(month)), localName)
	sourcepath := filepath.Join(config.Cfg.Basedir, relativePath)
	// 文件所在文件夹目录
	dirname := filepath.Dir(sourcepath)
	err = os.MkdirAll(dirname, 0666)
	if err != nil {
		result.SetData(err)
		c.JSON(200, result.SetResult(common.ResultMap["serverError"], ""))
		return
	}

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(f, sourcepath)
	if err != nil {
		c.JSON(200, result.Fail("上传失败", err))
		return
	}

	file := model.File{
		Name:      f.Filename,
		LocalName: localName,
		Uid:       uid,
		Ext:       ext[1:],
		Type:      filetype,
		MimeType:  mtype,
		Path:      relativePath,
		Size:      uint(f.Size),
	}
	res := db.Create(&file)
	if res.Error != nil {
		c.JSON(200, result.Fail("上传失败", res.Error))
		return
	}
	c.JSON(200, result.Success("", gin.H{
		"id": file.ID, "uid": file.Uid, "ext": file.Ext, "name": file.Name, "size": file.Size,
	}))
}

// Download 下载文件
func Download(c *gin.Context) {
	result := &common.Result{}
	id := c.Param("id")
	uid := util.Basename(id)
	db := c.Value("DB").(*gorm.DB)
	var file model.File

	res := db.First(&file, "uid", uid)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, result.SetResult(common.ResultMap["notFound"], ""))
		return
	}
	ext := filepath.Ext(id)
	if len(ext) > 0 && file.Ext != ext[1:] {
		c.JSON(http.StatusNotFound, result.SetResult(common.ResultMap["notFound"], ""))
		return
	}
	sourcePath := filepath.Join(config.Cfg.Basedir, file.Path)
	if !util.CheckFileIsExist(sourcePath) {
		c.JSON(http.StatusNotFound, result.SetResult(common.ResultMap["notFound"], ""))
		return
	}
	c.File(sourcePath)
}

// DownloadFromUrl 下载 url返回的数据，可下载第三方媒体文件
func DownloadFromUrl(c *gin.Context) {
	result := &common.Result{}
	url := c.Query("url")

	res, err := http.Get(url)
	if err != nil {
		c.JSON(200, result.Fail("", err))
		return
	}
	defer res.Body.Close()
	io.Copy(c.Writer, res.Body)
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
