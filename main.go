package main

import (
	// "context"

	"fmt"
	"gin_app/app/middleware"
	"gin_app/app/router"
	"gin_app/app/service"
	"gin_app/app/util"
	"gin_app/config"
	"io/ioutil"
	"reflect"
	"regexp"

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
	readRouters(router)

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

// 读取router下的路由组
func readRouters(g *gin.RouterGroup) {
	var funcNames []string
	fileInfo, _ := ioutil.ReadDir("app/router")
	if len(fileInfo) == 0 {
		return
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		// 匹配 .go结尾的文件
		re := regexp.MustCompile(`^\w+\.go$`)
		if re.MatchString(name) {
			basename := util.Basename(name)
			basename = util.UpperFirst(basename)
			funcNames = append(funcNames, basename)
		}
	}
	// 获取反射值
	value := reflect.ValueOf(&router.Router{})
	in := []reflect.Value{reflect.ValueOf(g)}
	for _, name := range funcNames {
		fn := value.MethodByName(name) //通过反射获取它对应的函数
		if fn.IsNil() {
			continue
		}
		fn.Call(in)
	}
}
