package router

import (
	"gin_app/app/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) File(g *gin.RouterGroup) {
	rg := g.Group("/file")
	{
		rg.GET("/:id", service.Download)
		rg.POST("/upload", service.Upload)
		rg.GET("/thumb/:id", service.DownloadThumb)
		rg.GET("/thumbInfo", service.ThumbInfo)
		rg.GET("/downloadFromUrl", service.DownloadFromUrl)
		rg.POST("/word", service.ExtractWord)
		rg.GET("/towebp", service.ConvertToWebp)
	}
}
