package main

import (
	"gin_app/api"
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

	router.Run() // listen and serve on 0.0.0.0:8080
}
