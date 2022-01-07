package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin_app/app/model"
	"gin_app/app/util"
	"gin_app/config"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
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
	log := c.Value("Logger").(*logrus.Entry)
	isSchedule := c.Query("schedule")

	res, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		fmt.Println("info err:", err)
		return
	}
	defer res.Body.Close()

	// 方法一：转成map对象后再转格式化的json对象
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// var result map[string]interface{}
	// if err := json.Unmarshal([]byte(body), &result); err != nil {
	// 	fmt.Println("err:", err)
	// }
	// imgURL := result["images"].(map[string]interface{})
	// fmt.Println(result, imgURL)

	// 转为json编码字符串并格式化输出
	// formatRes, err := json.MarshalIndent(result, "", "  ") //这里返回的data值，类型是[]byte
	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// }
	// fmt.Println(string(formatRes))

	// 方法二：解析为json对象
	bingRes := BingRes{}
	json.NewDecoder(res.Body).Decode(&bingRes)
	// 打印返回信息
	fmt.Println("bingRes:", bingRes)

	imgInfo := bingRes.Images[0]
	imgURL := util.WriteString("https://cn.bing.com", imgInfo.URL)
	res1, err := http.Get(imgURL)
	if err != nil {
		fmt.Println("img err:", err)
		return
	}
	defer res1.Body.Close()

	// Body 是 ReadCloser,只能读一次,不能 Seek ,只能把 Body 读出来, 保存到 buffer里面
	imgByte, err1 := ioutil.ReadAll(res1.Body)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	// res1.Body = ioutil.NopCloser(bytes.NewReader(imgByte))
	imgReader := bytes.NewReader(imgByte)

	// 使用固定的32K缓冲区，因此无论源数据多大，都只会占用32K内存空间
	io.Copy(c.Writer, imgReader)

	if isSchedule != "1" {
		return
	}

	db := c.Value("DB").(*gorm.DB)
	// 判断当天壁纸是否已下载
	var bing model.Bing
	res2 := db.Where("created_at >= ?", time.Now().Format("2006-01-02")).Limit(1).Find(&bing)
	if res2.Error != nil {
		log.Errorln(res2.Error)
		return
	}
	if res2.RowsAffected > 0 {
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
	result := uploadResult{}
	err = json.NewDecoder(fileRes.Body).Decode(&result)
	if err != nil {
		log.Errorln(err)
		return
	}
	if result.Code != 200 {
		log.Errorf("file upload failed: %v", result)
		return
	}
	releaseAt, _ := time.Parse("20060102", imgInfo.Enddate)
	bing = model.Bing{
		FileId:    result.Data.ID,
		Url:       imgURL,
		Hsh:       imgInfo.Hsh,
		Desc:      imgInfo.Copyright,
		ReleaseAt: releaseAt,
	}
	db.Create(&bing)

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
	db := c.Value("DB").(*gorm.DB)

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var bing []model.Bing
	var count int64
	tx := db

	tx.Joins("left join file on `file`.id = `bing`.file_id and `file`.is_del = 0").Model(&model.Bing{}).Count(&count)

	tx = tx.Preload("File", func(db *gorm.DB) *gorm.DB {
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

	c.JSON(200, gin.H{
		"count": count,
		"data":  bing,
	})
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
