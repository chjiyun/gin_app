package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util"
	"gin_app/config"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
type uploadResult struct {
	Code int
	Msg  string
	Data model.File
}

// GetImg 获取远程图片并返回
func GetImg(c *gin.Context) {
	isSchedule := c.Query("schedule")
	isUHD := c.Query("uhd")
	log := c.Value("Logger").(*logrus.Entry)
	db := c.Value("DB").(*gorm.DB)
	var bing model.Bing
	var file model.File

	res, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		fmt.Println("info err:", err)
		return
	}
	defer res.Body.Close()

	// 方法一：转成map对象后再转格式化的json对象

	// 方法二：解析为json对象
	bingRes := BingRes{}
	json.NewDecoder(res.Body).Decode(&bingRes)
	// 打印返回信息
	fmt.Println("bingRes:", bingRes)
	imgInfo := bingRes.Images[0]

	// 非定时任务请求时 检查本地文件是否已下载
	res2 := db.Where("created_at >= ?", time.Now().Format("2006-01-02")).Limit(1).Find(&bing)
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
		fmt.Println("img err:", err)
		return
	}
	defer res1.Body.Close()

	// Body 是 ReadCloser,只能读一次,不能 Seek ,只能把 Body 读出来, 保存到 buffer里面
	imgByte, err1 := io.ReadAll(res1.Body)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	imgReader := bytes.NewReader(imgByte)

	// 使用固定的32K缓冲区，因此无论源数据多大，都只会占用32K内存空间
	io.Copy(c.Writer, imgReader)

	if isSchedule != "1" {
		return
	}

	fileName := time.Now().Format("2006-01-02") + "." + imgInfo.Hsh[:16] + ".jpg"
	// sourcePath := filepath.Join("files", fileName)

	fd := map[string]interface{}{
		"file":     &imgByte,
		"filename": fileName,
	}
	uploadUrl := util.WriteString("http://127.0.0.1:", config.Cfg.Server.Port, "/api/file/upload")
	fileRes, err := util.SendFormData(uploadUrl, "file", fd)
	if err != nil {
		log.Errorf("error in SendFormData: %v", err)
		return
	}
	defer fileRes.Body.Close()
	upResult := uploadResult{}
	err = json.NewDecoder(fileRes.Body).Decode(&upResult)
	if err != nil {
		log.Errorln(err)
		return
	}
	if upResult.Code != 0 {
		log.Errorf("file upload failed: %v", upResult)
		return
	}
	releaseAt, _ := time.Parse("20060102", imgInfo.Enddate)
	bing = model.Bing{
		FileId:    upResult.Data.ID,
		Url:       imgURL,
		Hsh:       imgInfo.Hsh,
		Desc:      imgInfo.Copyright,
		ReleaseAt: releaseAt,
	}
	db.Create(&bing)
	db.Model(&upResult.Data).Update("desc", imgInfo.Copyright)

	// var f *os.File
	// defer f.Close()
	// if util.CheckFileIsExist(sourcePath) { //如果文件存在
	// 	// f, err1 = os.OpenFile(sourcePath, os.O_APPEND, 0666) //打开文件
	// 	fmt.Println("文件已存在")
	// 	return
	// } else {
	// 	f, err1 = os.Create(sourcePath) //创建文件
	// }
	// if err1 != nil {
	// 	panic(err1)
	// }
	// writer := bufio.NewWriter(f) //创建新的 Writer 对象

	// if err1 != nil {
	// 	fmt.Println(err1)
	// }
	// n, _ := writer.Write(imgByte)
	// fmt.Printf("写入 %d 个字节\n", n)
	// writer.Flush()

}

// GetAllBing 搜索符合条件的记录
func GetAllBing(c *gin.Context) {
	r := result.New()
	db := c.Value("DB").(*gorm.DB)

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var bing []model.Bing
	var count int64

	db.Model(&model.Bing{}).Count(&count)

	tx := db.Preload("File", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, uid, ext, name, size")
	}).Omit("url", "hsh", "updated_at")

	if startTime != "" {
		tx = tx.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		tx = tx.Where("created_at < ?", endTime)
	}
	if page > 0 && pageSize > 0 {
		tx = tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	tx.Order("created_at desc").Find(&bing)

	r.SetData(gin.H{
		"count": count,
		"rows":  bing,
	})
	c.JSON(200, r)
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
