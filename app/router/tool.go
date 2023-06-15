package router

import (
	"gin_app/app/controller/toolController"
	"github.com/gin-gonic/gin"
)

func (r Router) Tool(g *gin.RouterGroup) {
	rg := g.Group("/tool")
	{
		rg.GET("/ip", toolController.GetIpInfo)
	}
}
