package toolController

import (
	"gin_app/app/service/toolService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIpInfo(c *gin.Context) {
	r := toolService.GetIpInfo(c)
	c.JSON(http.StatusOK, r)
}
