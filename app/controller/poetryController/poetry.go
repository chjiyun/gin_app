package poetryController

import (
	"gin_app/app/controller/poetryController/poetryVo"
	"gin_app/app/result"
	"gin_app/app/service/poetryService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchPoetry(c *gin.Context) {
	r := result.New()
	var reqVo poetryVo.PoetrySearchReqVo
	if err := c.ShouldBindQuery(&reqVo); err != nil {
		c.JSON(http.StatusOK, r.Fail("invalid query params"))
		return
	}
	res, err := poetryService.SearchPoetry(c, reqVo)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func GetPoetry(c *gin.Context) {
	r := result.New()
	id := c.Param("id")
	res := poetryService.GetPoetry(c, id)
	c.JSON(http.StatusOK, r.Success(res))
}

func CreatePoetry(c *gin.Context) {
	r := result.New()
	res, err := poetryService.CreatePoetry(c)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}

func UpdatePoetry(c *gin.Context) {
	r := result.New()

	c.JSON(http.StatusOK, r)
}

func ImportPoetry(c *gin.Context) {
	r := result.New()
	res, err := poetryService.PoetryImport(c)
	if err != nil {
		c.JSON(http.StatusOK, r.FailErr(err))
		return
	}
	c.JSON(http.StatusOK, r.Success(res))
}
