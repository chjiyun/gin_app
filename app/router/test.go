package router

import (
	"gin_app/app/test"

	"github.com/gin-gonic/gin"
)

func (r *Router) Test(g *gin.RouterGroup) {
	rg := g.Group("/test")
	{
		rg.GET("/index", test.For)
		rg.GET("/map", test.Map)
		rg.GET("/arr", test.Arr)
		rg.GET("/json", test.Json)
		rg.GET("/str", test.String)
		rg.GET("/int", test.Int)
		rg.GET("/snowflake", test.Snowflake)
		rg.GET("/chan", test.Channel)
	}
}
