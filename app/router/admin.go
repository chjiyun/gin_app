package router

import (
	"gin_app/app/controller/adminController"
	"github.com/gin-gonic/gin"
)

func (r Router) Admin(g *gin.RouterGroup) {
	rg := g.Group("/admin")
	{
		rg.GET("/wallpaper", adminController.GetWallPaperPage)
		rg.POST("/wallpaper/audit", adminController.AuditWallPaper)
	}
}
