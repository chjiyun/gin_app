package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func getBingImg() {
	// 新建一个定时任务对象
	// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
	//crontab := cron.New()  默认从分开始进行时间调度
	crontab := cron.New(cron.WithSeconds()) //精确到秒
	//定义定时器调用的任务函数
	task := func() {
		fmt.Println("hello world", time.Now())
	}
	//定时任务
	spec := "*/5 * * * * ?" //cron表达式，每五秒一次
	// spec := "0 5 0 * * ?" //凌晨0点5分

	crontab.AddFunc(spec, task)
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {} //阻塞主线程停止
}
