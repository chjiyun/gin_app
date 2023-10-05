package bingController

import (
	"gin_app/app/controller/bingController/bingVo"
	"gin_app/app/result"
	"gin_app/app/service/bingService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWallPaper(c *gin.Context) {
	r := result.New()
	var reqVo bingVo.WallPaperReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := bingService.GetWallPaper(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func GetAllBing(c *gin.Context) {
	r := result.New()
	var reqVo bingVo.BingPageReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := bingService.GetAllBing(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func GetImgFromBing(c *gin.Context) {

}
