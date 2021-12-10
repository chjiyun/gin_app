package main

import (
	// "context"

	"fmt"
	"gin_app/app"
	"gin_app/app/middleware"
	"gin_app/app/service"
	"gin_app/config"

	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
)

func main() {
	// 初始化配置
	config.Init()

	// r := gin.Default()
	r := gin.New()
	r.Use(middleware.LoggerToFile(), middleware.SetContext(), gin.Recovery())

	// 简单的路由组: api

	r.GET("/", service.Index)
	router := r.Group("/api")
	app.ReadRouters(router)

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

	app.InitSchedule()
	fmt.Println("schedule init success...")

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	// r.Run(":8000") for a hard coded port
	r.Run(":" + config.Cfg.Server.Port)
}
