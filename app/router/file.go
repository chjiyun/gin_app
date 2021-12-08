package router

import (
	"gin_app/app/api"

	"github.com/gin-gonic/gin"
)

func (r *Router) File(g *gin.RouterGroup) {
	rg := g.Group("/file")
	{
		rg.GET("/:id", api.Download)
		rg.POST("/upload", api.Upload)
		rg.POST("/word", api.ExtractWord)
	}
}
