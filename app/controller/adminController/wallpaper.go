package adminController

import (
	"gin_app/app/controller/adminController/adminVo"
	"gin_app/app/result"
	"gin_app/app/service/adminService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWallPaperPage(c *gin.Context) {
	r := result.New()
	var reqVo adminVo.WallPaperReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := adminService.GetWallpaperPage(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func AuditWallPaper(c *gin.Context) {
	r := result.New()
	var reqVo adminVo.WallPaperAuditReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := adminService.AuditWallPaper(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}
