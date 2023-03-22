package toolController

import (
	"gin_app/app/result"
	"gin_app/app/service/toolService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIpInfo(c *gin.Context) {
	r := result.New()
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusOK, r.Fail("ip is required"))
		return
	}
	res, err := toolService.GetIpInfo(c, ip)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}
