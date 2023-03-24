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
