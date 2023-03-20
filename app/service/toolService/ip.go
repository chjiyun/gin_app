package toolService

import (
	"gin_app/app/common"
	"gin_app/app/common/myError"
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

func GetIpInfo(c *gin.Context, ip string) (*IPInfo, error) {
	client := req.C()
	var ipInfo IPInfo
	resp, err := client.R().
		SetQueryParam("ip", ip).
		SetSuccessResult(&ipInfo).
		Get(common.SEARCH_IP_URL)

	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, myError.New("远程请求错误")
	}
	return &ipInfo, nil
}
