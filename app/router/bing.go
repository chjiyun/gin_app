package router

import (
	"gin_app/app/controller/bingController"
	"gin_app/app/service/bingService"

	"github.com/gin-gonic/gin"
)

func (r Router) Bing(g *gin.RouterGroup) {
	rg := g.Group("/bing")
	{
		rg.GET("/getImg", bingService.GetImg)
		rg.GET("/getAllBing", bingController.GetAllBing)
		rg.GET("/wallpaper", bingController.GetWallPaper)
		rg.GET("/zip", bingService.GetBingZip)
		rg.POST("/wallpaper", bingController.CreateWallPaper)
		rg.PUT("/wallpaper", bingController.UpdateWallPaper)
		rg.DELETE("/wallpaper/:id", bingController.DeleteWallPaper)
		rg.POST("/wallpaper/validate", bingController.ValidateWallPaper)
		rg.POST("/wallpaper/audit", bingController.AuditWallPaper)
	}
}
