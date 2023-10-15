package bingController

import (
	"gin_app/app/controller/bingController/bingVo"
	"gin_app/app/result"
	"gin_app/app/service/bingService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWallPaperPage(c *gin.Context) {
	r := result.New()
	var reqVo bingVo.WallPaperReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := bingService.GetWallPaperPage(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func GetWallPaper(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	res, err := bingService.GetWallPaper(c, id)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func ValidateWallPaper(c *gin.Context) {
	r := result.New()
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := bingService.ValidateWallPaper(c, f)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func CreateWallPaper(c *gin.Context) {
	r := result.New()
	var reqVo bingVo.WallPaperCreateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := bingService.AddWallPaper(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func UpdateWallPaper(c *gin.Context) {
	r := result.New()
	var reqVo bingVo.WallPaperUpdateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := bingService.UpdateWallPaper(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func DeleteWallPaper(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	res, err := bingService.DeleteWallPaper(c, id)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}
