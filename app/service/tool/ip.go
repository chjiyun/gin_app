package tool

import (
	"fmt"
	"gin_app/app/common"
	"gin_app/app/result"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

type IPInfo struct {
	Ip        string `json:"ip"`
	Net       string `json:"net"`
	Isp       string `json:"isp"`
	Country   string `json:"country"`
	ShortName string `json:"short_name"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Area      string `json:"area"`
	Code      int    `json:"code"`
	Desc      string `json:"desc"`
}

func GetIpInfo(c *gin.Context) {
	r := result.New()

	ip := c.Query("ip")
	if ip == "" {
		c.JSON(200, r.Fail("ip is not exist"))
		return
	}

	client := req.C()
	var ipInfo IPInfo
	resp, err := client.R().
		SetQueryParam("ip", ip).
		SetSuccessResult(&ipInfo).
		Get(common.SEARCH_IP_URL)

	if err != nil {
		fmt.Println(err)
		c.JSON(200, r.Fail("请求错误"))
		return
	}
	if !resp.IsSuccessState() {
		c.JSON(200, r.Fail(""))
		return
	}
	r.SetData(ipInfo)
	c.JSON(200, r)
}
