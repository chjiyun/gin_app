package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	x := "hello"
	res, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		fmt.Println("err1:", err)
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
	fmt.Println("bingRes:", bingRes)
	str := bingRes.Images[0].URL
	// 最高效的字符串拼接方式
	var build strings.Builder
	build.WriteString("https://cn.bing.com")
	build.WriteString(str)
	imgURL := build.String()
	fmt.Println(imgURL)

	res1, err := http.Get(imgURL)
	if err != nil {
		fmt.Println("err2:", err)
		return
	}
	defer res1.Body.Close()
	fmt.Println(res1.Body)

	c.JSON(http.StatusOK, x)
}
