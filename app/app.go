package app

import (
	"gin_app/app/router"
	"gin_app/app/schedule"
	"gin_app/app/util"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// 读取router下的路由组
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
		job := result[0].Interface().(schedule.GinSchedule)
		if job.Disable {
			continue
		}
		crontab.AddFunc(job.Cron, job.Task)
	}
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	// defer crontab.Stop()

	// 根据实际情况进行控制
	// select {} //阻塞主线程停止
}

// 利用反射动态调用 --- 测试
// func DynamicFunc() {
// 	funcMap := map[string]interface{}{
// 		"GetBingImg": schedule.GetBingImg,
// 	}
// 	for k := range funcMap {
// 		fmt.Println("start schedule:", k)
// 		result, err := Call(funcMap, k)
// 		if err != nil {
// 			continue
// 		}
// 		for _, v := range result {
// 			// 打印返回值和类型
// 			fmt.Printf("type=%v, value=%+v\n", v.Type(), v.Interface().(schedule.GinSchedule))
// 		}
// 	}
// }
