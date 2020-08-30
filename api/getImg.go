package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Arr 数组和切片
func GetImg(c *gin.Context) {
	x := "hello"
	res, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 转成map对象后再转格式化的json对象
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		fmt.Println("err:", err)
	}
	formatRes, err := json.MarshalIndent(result, "", "  ") //这里返回的data值，类型是[]byte
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println(string(formatRes))

	c.JSON(http.StatusOK, x)
}
