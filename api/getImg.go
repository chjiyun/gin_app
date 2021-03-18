package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gin_app/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
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
}

// GetImg 获取远程图片并返回
func GetImg(c *gin.Context) {
	// x := "hello"
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

	str := bingRes.Images[0].URL
	// 最高效的字符串拼接方式
	var build strings.Builder
	build.WriteString("https://cn.bing.com")
	build.WriteString(str)
	imgURL := build.String()

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
	res1.Body = ioutil.NopCloser(bytes.NewReader(imgByte))

	// 使用固定的32K缓冲区，因此无论源数据多大，都只会占用32K内存空间
	io.Copy(c.Writer, res1.Body)

	filename := filepath.Join("files", bingRes.Images[0].Hsh+".jpg")

	var f *os.File
	if util.CheckFileIsExist(filename) { //如果文件存在
		// f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
		return
	} else {
		f, err1 = os.Create(filename) //创建文件
	}
	defer f.Close()
	if err1 != nil {
		panic(err1)
	}
	writer := bufio.NewWriter(f) //创建新的 Writer 对象

	if err1 != nil {
		fmt.Println(err1)
	}
	n, _ := writer.Write(imgByte)
	fmt.Printf("写入 %d 个字节\n", n)
	writer.Flush()

}
