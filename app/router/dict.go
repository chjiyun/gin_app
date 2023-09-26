package router

import (
	"gin_app/app/controller/dictController"
	"github.com/gin-gonic/gin"
)

func (r Router) Dict(g *gin.RouterGroup) {
	rg := g.Group("/dict")
	{
		rg.GET("/type", dictController.GetDictType)
		rg.GET("/type/all", dictController.GetAllDictType)
		rg.POST("/type", dictController.CreateDictType)
		rg.PUT("/type", dictController.UpdateDictType)
		rg.DELETE("/type/:id", dictController.DeleteDictType)
		rg.GET("/value", dictController.GetDictValue)
		rg.GET("/value/list", dictController.GetDictValueByType)
		rg.POST("/value", dictController.CreateDictValue)
		rg.PUT("/value", dictController.UpdateDictValue)
		rg.DELETE("/value/:id", dictController.DeleteDictValue)
	}
}
