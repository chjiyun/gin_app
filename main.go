package main

import (
	// "context"

	"fmt"
	"gin_app/api"
	"gin_app/config"
	"gin_app/middleware"
	"gin_app/service"
	"gin_app/test"
	"gin_app/util"

	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
)

func main() {
	// 初始化配置
	config.Init()

	// r := gin.Default()
	// r.Use(middleware.SetDB())
	r := gin.New()
	r.Use(middleware.LoggerToFile(), middleware.SetContext(), gin.Recovery())

	// 简单的路由组: api
	v1 := r.Group("/api")
	{
		v1.GET("/index", api.Index)
		v1.GET("/user", api.Users)
		v1.GET("/img", api.GetImg)
		v1.GET("/company", service.GetCompanies)

		file := v1.Group("/file")
		file.GET("/:id", api.Download)
		file.POST("/upload", api.Upload)
		file.POST("/word", api.ExtractWord)
	}
	v2 := r.Group("/test")
	{
		v2.GET("/index", test.For)
		v2.GET("/map", test.Map)
		v2.GET("/arr", test.Arr)
		v2.GET("/json", test.Json)
		v2.GET("/str", test.String)
		v2.GET("/int", test.Int)
		v2.GET("/snowflake", test.Snowflake)
		v2.GET("/chan", test.Channel)
	}

	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: r,
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

	var options = idgen.NewIdGeneratorOptions(1)
	idgen.SetIdGenerator(options)
	fmt.Println("雪花算法生成器初始化完成>>>")

	util.InitSchedule()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	// r.Run(":8000") for a hard coded port
	r.Run(":" + config.Cfg.Server.Port)
}
