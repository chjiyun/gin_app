package router

import (
	"gin_app/app/controller/bingController"
	"gin_app/app/service"

	"github.com/gin-gonic/gin"
)

func (r Router) Bing(g *gin.RouterGroup) {
	rg := g.Group("/bing")
	{
		rg.GET("/getImg", service.GetImg)
		rg.GET("/getAllBing", bingController.GetAllBing)
		rg.GET("/wallpaper", bingController.GetWallPaper)
		rg.GET("/zip", service.GetBingZip)
	}
}
