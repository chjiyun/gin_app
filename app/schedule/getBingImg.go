package schedule

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// 任务都挂载在此结构体下面
type Schedule struct {
}

// 定时任务公共属性
type GinSchedule struct {
	Cron      string `json:"cron"`
	Immediate bool   `json:"immediate"`
	Disable   bool   `json:"disable"`
	Task      func()
}

// 定时获取必应的壁纸
func (s *Schedule) GetBingImg(c *context.Context) GinSchedule {
	task := GinSchedule{
		Cron:      "0 5 0 * * ?",
		Immediate: true,
	}
	// 任务执行内容...
	task.Task = func() {
		_, err := http.Get("http://127.0.0.1:8000/api/bing/getImg?schedule=1")
		if err != nil {
			fmt.Println("img err:", err)
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " 必应壁纸下载成功")
	}
	return task
}
