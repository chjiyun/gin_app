package service

import (
	"errors"
	"fmt"
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util"
	"gin_app/config"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/discord/lilliput"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": {lilliput.JpegQuality: 85},
	".png":  {lilliput.PngCompression: 7},
	".webp": {lilliput.WebpQuality: 85},
}

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

	// 关键：重置offset
	mFile.Seek(0, 0)

	imgSuffix := regexp.MustCompile(`jpg|jpeg|png$`)
	if imgSuffix.MatchString(ext) {
		width, height, err := getImageXY(mFile)
		if err != nil {
			return nil, myError.New("文件解码失败")
		}
		err = toWebp(c, file, width, height)
		if err != nil {
			return nil, myError.New("图片转webp失败")
		}
	}
	return &file, nil
}

// Download 下载文件
func Download(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	isThumb := c.Query("thumb")
	format := c.Query("format")
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

// ThumbInfo 获取image thumb
func ThumbInfo(c *gin.Context) {
	r := result.New()
	db := c.Value("DB").(*gorm.DB)
	uid := c.Query("uid")
	var file model.File

	// hasMany关系关联时表名要加s
	tx := db.Preload("Thumbs").Where("uid = ?", uid).First(&file)
	if tx.Error != nil {
		c.JSON(200, r.Fail("record not found"))
		return
	}
	r.SetData(file)
	c.JSON(200, r)
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

// getImageXY 获取图片宽高 px
func getImageXY(file io.Reader) (int, int, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}
	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y
	return width, height, nil
}

// toWebp 转webp格式
func toWebp(c *gin.Context, file model.File, width int, height int) error {
	db := c.Value("DB").(*gorm.DB)
	log := c.Value("Logger").(*zap.SugaredLogger)

	ext := ".webp"
	uid := idgen.NextId()
	localName := util.ToString(uid) + ext
	outputFilename := filepath.Join("files", "thumb", localName)
	err := transformImage(file.Path, outputFilename, width, height)
	if err != nil {
		log.Errorln("image transform failed", file.Path, err)
		return err
	}
	tInfo, err := os.Stat(outputFilename)
	if err != nil {
		log.Errorln("文件不存在", err)
		return err
	}
	name := util.Basename(file.Name) + ".thumb" + ext
	thumb := model.Thumb{
		Uid:       uid,
		FileId:    file.ID,
		Ext:       ext[1:],
		Name:      name,
		LocalName: localName,
		Path:      outputFilename,
		Size:      uint(tInfo.Size()),
		Width:     uint(width),
		Height:    uint(height),
	}
	res := db.Create(&thumb)
	if res.Error != nil {
		log.Errorln("thumb create error:", res.Error)
		return res.Error
	}
	return nil
}

func transformImage(inputFilename string, outputFilename string, outputWidth int, outputHeight int) error {
	inputBuf, err := os.ReadFile(inputFilename)
	if err != nil {
		return err
	}
	decoder, err := lilliput.NewDecoder(inputBuf)
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		return err
	}
	defer decoder.Close()
	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	if err != nil {
		return err
	}
	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	// create a buffer to store the output image, 10MB in this case
	outputImg := make([]byte, 5*1024*1024)

	// use user supplied filename to guess output type if provided
	// otherwise don't transcode (use existing type)
	outputType := filepath.Ext(outputFilename)
	if outputFilename == "" {
		outputType = filepath.Ext(inputFilename)
	}
	if outputWidth == 0 {
		outputWidth = header.Width()
	}
	if outputHeight == 0 {
		outputHeight = header.Height()
	}

	resizeMethod := lilliput.ImageOpsFit
	// if stretch {
	// 	resizeMethod = lilliput.ImageOpsResize
	// }
	if outputWidth == header.Width() && outputHeight == header.Height() {
		resizeMethod = lilliput.ImageOpsNoResize
	}

	opts := &lilliput.ImageOptions{
		FileType:             outputType,
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		EncodeOptions:        EncodeOptions[outputType],
	}
	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		return err
	}
	if outputFilename == "" {
		basename := util.Basename(inputFilename)
		str := util.WriteString(basename, ".", util.ToString(outputWidth), "_", util.ToString(outputHeight), outputType)
		outputFilename = str
	}
	if _, err := os.Stat(outputFilename); !os.IsNotExist(err) {
		str := fmt.Sprintf("%s 文件已存在", outputFilename)
		return errors.New(str)
	}
	err = os.WriteFile(outputFilename, outputImg, 0400)
	if err != nil {
		return err
	}
	return nil
}

// ConvertToWebp 图片转换成webp格式
func ConvertToWebp(c *gin.Context) {
	db := c.Value("DB").(*gorm.DB)
	log := c.Value("Logger").(*zap.SugaredLogger)
	var files []model.File

	db.Where("ext = ? or ext = ? or ext = ?", "jpg", "jpeg", "png").Find(&files)

	var errs []map[string]interface{}
	for _, file := range files {
		f, err := os.Open(file.Path)
		if err != nil {
			log.Errorln("文件打开失败", err)
			errs = append(errs, gin.H{
				"path":  file.Path,
				"error": err.Error(),
			})
			continue
		}
		width, height, err := getImageXY(f)
		if err != nil {
			log.Errorln("文件解码失败", err)
			errs = append(errs, gin.H{
				"path":  file.Path,
				"error": err.Error(),
			})
			continue
		}
		err = toWebp(c, file, width, height)
		if err != nil {
			errs = append(errs, gin.H{
				"path":  file.Path,
				"error": err.Error(),
			})
			continue
		}
	}

	c.JSON(200, errs)
}
