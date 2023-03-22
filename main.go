package main

import (
	"fmt"
	"gin_app/app"
	"gin_app/app/middleware"
	"gin_app/app/service"
	"gin_app/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
)

var (
	AppVersion = "Not provided"
	GoVersion  = "Not provided"
	BuildTime  = "Not provided"
	BuildUser  = "Not provided"
	CommitId   = "Not provided"
)

func main() {
	// 编译时注入信息
	args := os.Args
	if len(args) == 2 && (args[1] == "-v" || args[1] == "version") {
		fmt.Println("-------------------")
		fmt.Printf("App Version: %s\n", AppVersion)
		fmt.Printf("Go Version: %s\n", GoVersion)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Build User: %s\n", BuildUser)
		fmt.Printf("Git Commit Id: %s\n", CommitId)
		fmt.Println("-------------------")
		return
	}

	// 初始化配置
	config.Init()

	// r := gin.Default()
	r := gin.New()
	_ = r.SetTrustedProxies(nil)

	r.Use(middleware.LoggerToFile(), middleware.SetTimeout(), middleware.SetContext(), middleware.JWTAuth(), gin.Recovery())

	// 简单的路由组: api
	r.GET("/", service.Index)
	router := r.Group("/api")
	app.ReadRouters(router)

	app.RegisterValidation()

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
	fmt.Println(">>>雪花id生成器初始化完成")

	app.InitSchedule()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	// r.Run(":8000") for a hard coded port
	r.Run(":" + config.Cfg.Server.Port)
}
