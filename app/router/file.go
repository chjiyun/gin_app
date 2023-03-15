package router

import (
	"gin_app/app/controller"
	"gin_app/app/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) File(g *gin.RouterGroup) {
	rg := g.Group("/file")
	{
		rg.GET("/:id", service.Download)
		rg.POST("/upload", controller.Upload)
		rg.GET("/downloadFromUrl", service.DownloadFromUrl)
		rg.POST("/word", service.ExtractWord)
		rg.GET("/towebp", service.ConvertToWebp)
	}
}
