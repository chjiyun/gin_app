package fileController

import (
	"gin_app/app/result"
	"gin_app/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Upload(c *gin.Context) {
	r := result.New()
	res, err := service.Upload(c)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func DownloadThumb(c *gin.Context) {
	err := service.DownloadThumb(c)
	if err != nil {
		r := result.New()
		c.JSON(http.StatusNotFound, r.FailErr(err))
		return
	}
}
