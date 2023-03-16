package service

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/model"
	"gin_app/app/result"
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
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

// Upload 接收上传的文件
func Upload(c *gin.Context) (*model.File, error) {
	db := c.Value("DB").(*gorm.DB)

	f, err := c.FormFile("file")
	if err != nil {
		return nil, errors.New("上传失败")
	}

	ext := filepath.Ext(f.Filename)
	// mimetype := f.Header["Content-Type"][0]
	mfile, _ := f.Open()
	defer mfile.Close()
	mime, err := mimetype.DetectReader(mfile)
	if err != nil {
		return nil, errors.New("MIME type detect failed")
	}
	mtype := mime.String()
	var filetype string

	if i := strings.Index(mtype, "/"); i > 0 {
		filetype = mtype[:i]
	}
	uid := idgen.NextId()
	localName := util.ToString(uid) + ext
	year, month, _ := time.Now().Date()
	relativePath := filepath.Join("files", util.ToString(year), util.ToString(int(month)), localName)
	sourcepath := filepath.Join(config.Cfg.Basedir, relativePath)
	// 文件所在文件夹目录
	dirname := filepath.Dir(sourcepath)
	err = os.MkdirAll(dirname, 0666)
	if err != nil {
		return nil, errors.New("上传失败")
	}

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(f, sourcepath)
	if err != nil {
		return nil, errors.New("上传失败")
	}
	if true {
		return &model.File{Name: "test"}, nil
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
		return nil, errors.New("上传失败")
	}
	return &file, nil
}

// Download 下载文件
func Download(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	uid := util.Basename(id)
	ext := filepath.Ext(id)
	if ext == "" {
		c.JSON(http.StatusNotFound, r.FailType(common.FileNotFound))
		return
	}
	db := c.Value("DB").(*gorm.DB)
	var file model.File
	ext = ext[1:]

	res := db.Where("uid = ? AND ext = ?", uid, ext).First(&file)
	if res.Error != nil {
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
