package fileController

import (
	"gin_app/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Upload(c *gin.Context) {
	r := service.Upload(c)
	c.JSON(http.StatusOK, r)
}
