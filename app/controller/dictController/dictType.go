package dictController

import (
	"gin_app/app/controller/dictController/dictVo"
	"gin_app/app/result"
	"gin_app/app/service/dictService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDictType(c *gin.Context) {
	r := result.New()
	keyword := c.Query("keyword")
	res, err := dictService.GetDictType(c, keyword)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func GetAllDictType(c *gin.Context) {
	r := result.New()
	res, err := dictService.GetAllDictType(c)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func CreateDictType(c *gin.Context) {
	r := result.New()
	var reqVo dictVo.DictTypeCreateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := dictService.CreateDictType(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func UpdateDictType(c *gin.Context) {
	r := result.New()
	var reqVo dictVo.DictTypeUpdateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	res, err := dictService.UpdateDictType(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func DeleteDictType(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	res, err := dictService.DeleteDictType(c, id)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}
