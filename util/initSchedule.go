package util

import (
	"fmt"
	"gin_app/schedule"

	"github.com/robfig/cron/v3"
)

// InitSchedule 初始化当前目录所有定时任务
func InitSchedule() {
	funcMap := map[string]interface{}{
		"GetBingImg": schedule.GetBingImg,
	}
	// 新建一个定时任务对象
	// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
	// crontab := cron.New()  默认从分开始进行时间调度
	crontab := cron.New(cron.WithSeconds()) //精确到秒

	//定时任务
	// spec := "*/5 * * * * ?" //cron表达式，每五秒一次
	// spec := "0 5 0 * * ?" //凌晨0点5分
	for k := range funcMap {
		fmt.Println("start schedule:", k)
		result, err := Call(funcMap, k)
		if err != nil {
			continue
		}
		obj := result[0].Interface().(schedule.GinSchedule)
		crontab.AddFunc(obj.Cron, obj.Task)
	}
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	// defer crontab.Stop()

	// 根据实际情况进行控制
	// select {} //阻塞主线程停止
}

// 利用反射动态调用 --- 测试
func DynamicFunc() {
	funcMap := map[string]interface{}{
		"GetBingImg": schedule.GetBingImg,
	}
	for k := range funcMap {
		fmt.Println("start schedule:", k)
		result, err := Call(funcMap, k)
		if err != nil {
			continue
		}
		for _, v := range result {
			// 打印返回值和类型
			fmt.Printf("type=%v, value=%+v\n", v.Type(), v.Interface().(schedule.GinSchedule))
		}
	}
}
