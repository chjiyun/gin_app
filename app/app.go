package app

import (
	"fmt"
	"gin_app/app/router"
	"gin_app/app/schedule"
	"gin_app/app/util"
	"gin_app/app/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// ReadRouters 读取router下的路由组
func ReadRouters(g *gin.RouterGroup) {
	var funcNames = util.GetFileBasename("app/router", []string{"go"})
	if len(funcNames) == 0 {
		return
	}
	// 获取反射值
	value := reflect.ValueOf(&router.Router{})
	in := []reflect.Value{reflect.ValueOf(g)}
	for _, fnName := range funcNames {
		fn := value.MethodByName(fnName) //通过反射获取它对应的函数
		if fn.Kind() != reflect.Func || fn.IsNil() {
			continue
		}
		fn.Call(in)
	}
}

// InitSchedule 初始化定时任务配置，自动添加文件下所有任务到队列
func InitSchedule() {
	names := util.GetFileBasename("app/schedule", []string{"go"})
	if len(names) == 0 {
		return
	}
	// 新建一个定时任务对象
	// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
	// crontab := cron.New()  默认从分开始进行时间调度
	crontab := cron.New(cron.WithSeconds()) //精确到秒

	value := reflect.ValueOf(&schedule.Schedule{})
	in := []reflect.Value{}
	for _, fnName := range names {
		fn := value.MethodByName(fnName) //通过反射获取它对应的函数
		if fn.Kind() != reflect.Func || fn.IsNil() {
			continue
		}
		// 拿到定时任务配置结构体 GinSchedule
		result := fn.Call(in)
		job, ok := result[0].Interface().(schedule.GinSchedule)
		if !ok || job.Disable {
			continue
		}
		crontab.AddFunc(job.Cron, job.Task)
	}
	crontab.Start()
	fmt.Println(">>>schedule init successful")
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	// defer crontab.Stop()

	// 根据实际情况进行控制
	// select {} //阻塞主线程停止
}

func RegisterValidation() {
	vf := validation.ValidateFunc{}
	typ := reflect.TypeOf(vf)
	val := reflect.ValueOf(vf)
	if val.Kind() != reflect.Struct {
		return
	}
	// 获取到该结构体有多少个方法
	numOfMethod := val.NumMethod()
	if numOfMethod == 0 {
		return
	}
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	count := 0
	for i := 0; i < numOfMethod; i++ {
		fn, ok := val.Method(i).Interface().(func(fl validator.FieldLevel) bool)
		if !ok {
			continue
		}
		// 注册自定义校验函数
		err := validate.RegisterValidation(typ.Method(i).Name, fn)
		if err != nil {
			fmt.Println(typ.Method(i).Name, err)
			continue
		}
		count++
	}
	if count == numOfMethod {
		fmt.Println(">>>自定义校验函数注册完成")
	}
}
