package router

import (
	"gin_app/api"

	"github.com/gin-gonic/gin"
)

func (r *Router) Home(g *gin.RouterGroup) {
	rg := g.Group("/home")
	{
		rg.GET("/user", api.Users)
		rg.GET("/bing", api.GetImg)
	}
}
