package router

import (
	"gin_app/service"

	"github.com/gin-gonic/gin"
)

type Router struct {
}

func (r *Router) Company(g *gin.RouterGroup) {
	rg := g.Group("/company")
	{
		rg.GET("/list", service.GetCompanies)
	}
}
