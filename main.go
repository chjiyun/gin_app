package main

import (
	"gin_app/api"
	"gin_app/test"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/index", api.Index)
		v1.GET("/user", api.Users)
	}
	v2 := router.Group("/test")
	{
		v2.GET("/index", test.DataType)
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}
