package router

import (
	"gin_app/app/service/tool"

	"github.com/gin-gonic/gin"
)

func (r *Router) Tool(g *gin.RouterGroup) {
	rg := g.Group("/tool")
	{
		rg.GET("/ip", tool.GetIpInfo)
	}
}
