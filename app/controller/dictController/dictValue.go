package dictController

import (
	"gin_app/app/controller/dictController/dictVo"
	"gin_app/app/result"
	"gin_app/app/service/dictService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDictValue(c *gin.Context) {
	r := result.New()
	var reqVo dictVo.DictValueReqVo
	err := c.ShouldBindJSON(&reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res := dictService.GetDictValue(c, reqVo)
	c.JSON(http.StatusOK, r.Success(res))
}

func CreateDictValue(c *gin.Context) {
	r := result.New()
	var reqVo dictVo.DictValueCreateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := dictService.CreateDictValue(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func UpdateDictValue(c *gin.Context) {
	r := result.New()
	var reqVo dictVo.DictValueUpdateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := dictService.UpdateDictValue(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func DeleteDictValue(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	res, err := dictService.DeleteDictValue(c, id)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}
