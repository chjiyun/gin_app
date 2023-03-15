package controller

import (
	"gin_app/app/result"
	"gin_app/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Upload(c *gin.Context) {
	r := result.New()
	file, err := service.Upload(c)
	if err != nil {
		r.Fail(err.Error())
	}
	r.SetData(file)
	c.JSON(http.StatusOK, r)
}
