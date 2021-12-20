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
		rg.POST("/word", service.ExtractWord)
	}
}