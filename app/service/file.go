package service

import (
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util"
	"gin_app/config"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

// Upload 接收上传的文件
func Upload(c *gin.Context) (*model.File, error) {
	db := c.Value("DB").(*gorm.DB)

	f, err := c.FormFile("file")
	if err != nil {
		return nil, myError.NewET(common.FileNotFound)
	}

	ext := filepath.Ext(f.Filename)
	// mimetype := f.Header["Content-Type"][0]
	mFile, _ := f.Open()
	defer mFile.Close()
	mime, err := mimetype.DetectReader(mFile)
	if err != nil {
		return nil, err
	}
	mType := mime.String()
	var filetype string

	if i := strings.Index(mType, "/"); i > 0 {
		filetype = mType[:i]
	}
	uid := idgen.NextId()
	id := gonanoid.Must()
	localName := util.ToString(uid) + ext
	year, month, _ := time.Now().Date()
	relativePath := filepath.Join("files", util.ToString(year), util.ToString(int(month)), localName)
	sourcepath := filepath.Join(config.Cfg.Basedir, relativePath)
	// 文件所在文件夹目录
	dirname := filepath.Dir(sourcepath)
	err = os.MkdirAll(dirname, 0666)
	if err != nil {
		return nil, err
	}

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(f, sourcepath)
	if err != nil {
		return nil, myError.New("文件保存失败")
	}

	file := model.File{
		Name:      f.Filename,
		LocalName: localName,
		Uid:       uid,
		Ext:       ext[1:],
		Type:      filetype,
		MimeType:  mType,
		Path:      relativePath,
		Size:      uint(f.Size),
	}
	file.ID = id
	res := db.Create(&file)
	if res.Error != nil {
		return nil, myError.NewET(common.UnknownError)
	}
	return &file, nil
}

// Download 下载文件
func Download(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	realId := util.Basename(id)
	ext := filepath.Ext(id)
	db := c.Value("DB").(*gorm.DB)
	var file model.File
	if ext != "" {
		ext = ext[1:]
	}

	res := db.Where("id = ?", realId).First(&file)
	if res.Error != nil || file.Ext != ext {
		c.JSON(http.StatusNotFound, r.FailType(common.FileNotFound))
		return
	}
	sourcePath := filepath.Join(config.Cfg.Basedir, file.Path)
	if !util.CheckFileIsExist(sourcePath) {
		c.JSON(http.StatusNotFound, r.FailType(common.FileNotFound))
		return
	}
	c.File(sourcePath)
}

// DownloadFromUrl 下载 url返回的数据，可下载第三方媒体文件
func DownloadFromUrl(c *gin.Context) {
	r := result.New()
	url := c.Query("url")

	res, err := http.Get(url)
	if err != nil {
		c.JSON(200, r.Fail("请求异常"))
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

// convertToWebp 图片转换成webp格式
func ConvertToWebp(c *gin.Context) {
	db := c.Value("DB").(*gorm.DB)
	var files []model.File

	db.Find(&files)

	// for _, img := range files {

	// }

	c.JSON(200, files)
}
