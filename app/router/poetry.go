package router

import (
	"gin_app/app/controller/poetryController"
	"github.com/gin-gonic/gin"
)

func (r Router) Poetry(g *gin.RouterGroup) {
	rg := g.Group("/poetry")
	{
		rg.GET("/search", poetryController.SearchPoetry)
		rg.GET("/:id", poetryController.GetPoetry)
		rg.POST("/", poetryController.CreatePoetry)
		rg.PUT("/", poetryController.UpdatePoetry)
		rg.POST("/import", poetryController.ImportPoetry)
	}
}
