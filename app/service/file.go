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
	file.ID = id + ext
	res := db.Create(&file)
	if res.Error != nil {
		return nil, myError.NewET(common.UnknownError)
	}
	return &file, nil
}

// Save 临时文件保存到files目录
func Save(c *gin.Context, f *os.File) (string, error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return "", err
	}
	size := fileInfo.Size()
	name := fileInfo.Name()
	ext := filepath.Ext(name)
	// 必须重置偏移量
	if _, err = f.Seek(0, 0); err != nil {
		return "", err
	}
	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return "", err
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
	sourcePath := filepath.Join(config.Cfg.Basedir, relativePath)
	// 文件所在文件夹目录
	dirname := filepath.Dir(sourcePath)
	if err = os.MkdirAll(dirname, 0666); err != nil {
		return "", err
	}
	// 复制到另一个目录
	out, err := os.Create(sourcePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = f.Seek(0, 0); err != nil {
		return "", err
	}
	if _, err = io.Copy(out, f); err != nil {
		return "", err
	}
	db := c.Value("DB").(*gorm.DB)
	file := model.File{
		Name:      name,
		LocalName: localName,
		Uid:       uid,
		Ext:       ext[1:],
		Type:      filetype,
		MimeType:  mType,
		Path:      relativePath,
		Size:      uint(size),
	}
	file.ID = id + ext
	if err = db.Create(&file).Error; err != nil {
		return "", myError.NewET(common.UnknownError)
	}
	return file.ID, nil
}

// Download 下载文件
func Download(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	isThumb := c.Query("thumb")
	format := c.Query("format")
	ext := filepath.Ext(id)
	db := c.Value("DB").(*gorm.DB)
	var file model.File
	if ext != "" {
		ext = ext[1:]
	}

	res := db.Where("id = ?", id).First(&file)
	if res.Error != nil || file.Ext != ext {
		c.JSON(http.StatusNotFound, r.FailType(common.FileNotFound))
		return
	}
	sourcePath := file.Path
	// 返回thumb文件
	if isThumb != "" {
		var thumb model.Thumb
		tx := db.Where("file_id = ?", file.ID)
		if format != "" {
			tx = tx.Where("ext = ?", format)
		}
		tx.Select("id", "path").First(&thumb)
		if tx.Error != nil {
			c.JSON(http.StatusNotFound, r.Fail("thumb not found"))
			return
		}
		sourcePath = thumb.Path
	}
	sourcePath = filepath.Join(config.Cfg.Basedir, sourcePath)
	if !util.CheckFileIsExist(sourcePath) {
		c.JSON(http.StatusNotFound, r.FailType(common.FileNotFound))
		return
	}
	c.File(sourcePath)
}

func DownloadThumb(c *gin.Context) error {
	id := c.Param("id")
	ext := util.GetFileExt(id)
	db := c.Value("DB").(*gorm.DB)
	var thumb model.Thumb

	res := db.Where("id = ?", id).First(&thumb)
	if res.Error != nil || thumb.Ext != ext {
		return myError.NewET(common.FileNotFound)
	}
	sourcePath := filepath.Join(config.Cfg.Basedir, thumb.Path)
	if !util.CheckFileIsExist(sourcePath) {
		return myError.NewET(common.FileNotFound)
	}
	c.File(sourcePath)
	return nil
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
