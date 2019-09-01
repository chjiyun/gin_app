package main

import (
	"./services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/index", index.Index)
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}
