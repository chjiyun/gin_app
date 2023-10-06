package bingService

import (
	"bytes"
	"gin_app/app/common"
	"gin_app/app/controller/bingController/bingVo"
	"gin_app/app/model"
	"gin_app/app/service"
	"gin_app/app/util"
	"gin_app/config"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BingRes 接收接口响应
type BingRes struct {
	Images []ImgInfo `json:"images"`
}

// ImgInfo 图片详细信息
type ImgInfo struct {
	URL       string `json:"url"`
	Urlbase   string `json:"urlbase"`
	Copyright string `json:"copyright"`
	Hsh       string
	Enddate   string
}

// GetImg 获取远程图片并返回
func GetImg(c *gin.Context) {
	isSchedule := c.Query("schedule")
	isUHD := c.Query("uhd")
	log := c.Value("Logger").(*zap.SugaredLogger)
	db := c.Value("DB").(*gorm.DB)

	var bing model.Bing
	var file model.File

	res, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	// 方法一：转成map对象后再转格式化的json对象

	// 方法二：解析为json对象
	bingRes := BingRes{}
	if err = sonic.ConfigDefault.NewDecoder(res.Body).Decode(&bingRes); err != nil {
		log.Error(err)
		return
	}

	imgInfo := bingRes.Images[0]

	// 非定时任务请求时 检查本地文件是否已下载
	today := time.Now().Format(time.DateOnly)
	res2 := db.Where("created_at >= ? and is_bing=true", today).Limit(1).Find(&bing)
	if res2.Error != nil {
		log.Errorln(res2.Error)
		return
	}
	if isUHD == "" && res2.RowsAffected > 0 && config.Cfg.Env != gin.DebugMode {
		if isSchedule != "1" {
			fileRes := db.First(&file, "id", bing.FileId)
			if fileRes.Error != nil {
				log.Errorln(fileRes.Error)
				return
			}
			sourcePath := filepath.Join(config.Cfg.Basedir, file.Path)
			c.File(sourcePath)
		}
		return
	}

	if isUHD != "" {
		imgInfo.URL = strings.ReplaceAll(imgInfo.URL, "1920x1080", "UHD")
	}
	// 获取图片
	imgURL := util.WriteString("https://cn.bing.com", imgInfo.URL)
	res1, err := http.Get(imgURL)
	if err != nil {
		log.Error(err)
		return
	}
	defer res1.Body.Close()

	// Body 是 ReadCloser,只能读一次,不能 Seek ,只能把 Body 读出来, 保存到 buffer里面
	imgByte, err := io.ReadAll(res1.Body)
	if err != nil {
		log.Error(err)
		return
	}
	imgReader := bytes.NewReader(imgByte)

	// 使用固定的32K缓冲区，因此无论源数据多大，都只会占用32K内存空间
	io.Copy(c.Writer, imgReader)

	if isSchedule != "1" {
		return
	}

	fileName := time.Now().Format("2006-01-02") + "." + imgInfo.Hsh[:16] + ".jpg"
	sourcePath := filepath.Join(config.Cfg.Basedir, "files/temp", fileName)
	if util.CheckFileIsExist(sourcePath) {
		_ = os.Remove(sourcePath)
	}
	out, err := os.Create(sourcePath)
	if err != nil {
		return
	}
	if _, err = imgReader.Seek(0, 0); err != nil {
		return
	}
	if _, err = io.Copy(out, imgReader); err != nil {
		return
	}
	fileId, err := service.Save(c, out)
	if err != nil {
		log.Error(err)
		return
	}
	releaseAt, _ := time.Parse("20060102", imgInfo.Enddate)
	id := idgen.NextId()
	bing = model.Bing{
		FileId:    fileId,
		Url:       imgURL,
		Hsh:       imgInfo.Hsh,
		Desc:      imgInfo.Copyright,
		ReleaseAt: releaseAt,
		IsBing:    true,
	}
	bing.ID = id
	db.Create(&bing)
	out.Close()
	_ = os.Remove(sourcePath)
}

// GetAllBing 获取bing数据
func GetAllBing(c *gin.Context, reqVo bingVo.BingPageReqVo) (common.PageRes, error) {
	db := c.Value("DB").(*gorm.DB)

	var bing []model.Bing
	var count int64
	var respVo []bingVo.BingRespVo
	var pageRes common.PageRes

	tx := db.Model(&model.Bing{}).Where("is_bing = ?", true)

	if !reqVo.StartTime.IsZero() {
		tx = tx.Where("created_at >= ?", reqVo.StartTime)
	}
	if !reqVo.EndTime.IsZero() {
		tx = tx.Where("created_at < ?", reqVo.EndTime)
	}
	tx.Count(&count)
	tx = tx.Preload("File", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, uid, ext, name, size")
	}).Omit("url", "hsh", "pass", "updated_at", "is_del").
		Order("created_at desc").
		Limit(reqVo.PageSize).Offset((reqVo.Page - 1) * reqVo.PageSize).Find(&bing)

	_ = copier.Copy(&respVo, &bing)
	pageRes.Count = count
	pageRes.Rows = respVo
	return pageRes, nil
}

// GetBingZip 压缩下载bing图片
func GetBingZip(c *gin.Context) {
	db := c.Value("DB").(*gorm.DB)
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	tx := db
	var bing []model.Bing
	// var file []model.File

	if startTime != "" {
		tx = tx.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		tx = tx.Where("created_at < ?", endTime)
	}
	tx.Preload("File", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, path, name")
	}).Select("`bing`.id, `bing`.file_id").Find(&bing)

	// tx.Joins("File").Select("`bing`.id, `bing`.file_id").Find(&bing)

	// fmt.Println(file)

	filenames := make([]string, 0, len(bing))
	dst := make([]string, 0, len(bing))
	zipName := "bing_wallpaper.zip"
	for _, f := range bing {
		if f.File.Path != "" {
			path := filepath.Join(config.Cfg.Basedir, f.File.Path)
			dstpath := filepath.Join("bing_wallpaper", f.File.Name)
			filenames = append(filenames, path)
			dst = append(dst, dstpath)
		}
	}
	// header 在写入writer前设置
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+zipName)
	c.Header("Content-Transfer-Encoding", "binary")

	util.ZipFiles(&c.Writer, filenames, dst)
}
