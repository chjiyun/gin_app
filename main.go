package main

import (
	// "context"
	"gin_app/api"
	"gin_app/test"

	// "log"
	// "net/http"
	// "os"
	// "os/signal"
	// "time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 简单的路由组: api
	v1 := router.Group("/api")
	{
		v1.GET("/index", api.Index)
		v1.GET("/user", api.Users)
	}
	v2 := router.Group("/test")
	{
		v2.GET("/index", test.DataType)
		v2.GET("/map", test.Map)
	}

	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: router,
	// }

	// go func() {
	// 	// 服务连接
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()

	// // 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	// quit := make(chan os.Signal)
	// signal.Notify(quit, os.Interrupt)
	// <-quit
	// // log.Println("Shutdown Server ...")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown:", err)
	// }
	// log.Println("Server exiting")

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	// router.Run(":3000") for a hard coded port
	router.Run(":3000")
}
