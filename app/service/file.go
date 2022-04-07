package service

import (
	"errors"
	"fmt"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util"
	"gin_app/config"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/discord/lilliput"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
	"github.com/sirupsen/logrus"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": {lilliput.JpegQuality: 85},
	".png":  {lilliput.PngCompression: 7},
	".webp": {lilliput.WebpQuality: 85},
}

// Upload 接收上传的文件
func Upload(c *gin.Context) {
	r := result.New()
	// db := c.Value("DB").(*gorm.DB)

	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, r.Fail("上传失败", err))
		return
	}

	ext := filepath.Ext(f.Filename)
	// mimetype := f.Header["Content-Type"][0]
	mfile, _ := f.Open()
	defer mfile.Close()
	mime, err := mimetype.DetectReader(mfile)
	if err != nil {
		c.JSON(200, r.Fail("MIME type detect failed", err))
		return
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
		r.SetResult(result.ResultMap["serverError"], "").SetError(err)
		c.JSON(200, r)
		return
	}

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(f, sourcepath)
	if err != nil {
		c.JSON(200, r.Fail("上传失败", err))
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
	// res := db.Create(&file)
	// if res.Error != nil {
	// 	c.JSON(200, r.Fail("", res.Error))
	// 	return
	// }
	fmt.Println(file)

	// 关键：重置offset
	mfile.Seek(0, 0)
	width, height, err := getImageXY(mfile)
	if err != nil {
		c.JSON(200, r.Fail("文件解码失败", err))
		return
	}
	fmt.Println(width, height)

	r.SetData(gin.H{
		"id": file.ID, "uid": file.Uid, "ext": file.Ext, "name": file.Name, "size": file.Size,
	})
	c.JSON(200, r)
}

// Download 下载文件
func Download(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	uid := util.Basename(id)
	ext := filepath.Ext(id)
	if ext == "" {
		c.JSON(http.StatusNotFound, r.SetResult(result.ResultMap["notFound"], ""))
		return
	}
	db := c.Value("DB").(*gorm.DB)
	var file model.File
	ext = ext[1:]

	res := db.Where("uid = ? AND ext = ?", uid, ext).First(&file)
	if res.Error != nil {
		r.SetResult(result.ResultMap["notFound"], "").SetError(res.Error)
		c.JSON(http.StatusNotFound, r)
		return
	}
	sourcePath := filepath.Join(config.Cfg.Basedir, file.Path)
	if !util.CheckFileIsExist(sourcePath) {
		c.JSON(http.StatusNotFound, r.SetResult(result.ResultMap["notFound"], ""))
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
		c.JSON(200, r.Fail("", err))
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
func toWebp(c *gin.Context, file model.File) {
	db := c.Value("DB").(*gorm.DB)
	log := c.Value("Logger").(*logrus.Entry)

	ext := ".webp"
	uid := idgen.NextId()
	localName := util.ToString(uid) + ext
	outputFilename := filepath.Join("files", "thumb", localName)
	err := transformImage(file.Path, outputFilename, 0, 0)
	if err != nil {
		log.Errorln("image transform failed", file.Path, err)
		return
	}
	tInfo, err := os.Stat(outputFilename)
	if err != nil {
		log.Errorln("文件不存在", err)
		return
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
		Width:     0,
		Height:    0,
	}
	res := db.Create(&thumb)
	if res.Error != nil {
		log.Errorln("thumb create error:", res.Error)
	}
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
	outputImg := make([]byte, 10*1024*1024)

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

// convertToWebp 图片转换成webp格式
func ConvertToWebp(c *gin.Context) {
	db := c.Value("DB").(*gorm.DB)
	log := c.Value("Logger").(*logrus.Entry)
	var files []model.File

	db.Where("ext = ? or ext = ? or ext = ?", "jpg", "jpeg", "png").Find(&files)

	for _, img := range files {
		inputFilename := img.Path
		ext := ".webp"
		uid := idgen.NextId()
		localName := util.ToString(uid) + ext
		outputFilename := filepath.Join("files", "thumb", localName)
		err := transformImage(inputFilename, outputFilename, 0, 0)
		if err != nil {
			log.Errorln("image transform failed", inputFilename, err)
			continue
		}
		tInfo, err := os.Stat(outputFilename)
		if err != nil {
			log.Errorln("文件不存在", err)
			continue
		}
		name := util.Basename(img.Name) + "thumb" + ext
		thumb := model.Thumb{
			Uid:       uid,
			FileId:    img.ID,
			Ext:       ext[1:],
			Name:      name,
			LocalName: localName,
			Path:      outputFilename,
			Size:      uint(tInfo.Size()),
			Width:     0,
			Height:    0,
		}
		res := db.Create(&thumb)
		if res.Error != nil {
			log.Errorln("thumb create error:", res.Error)
			continue
		}
	}

	c.JSON(200, files)
}
