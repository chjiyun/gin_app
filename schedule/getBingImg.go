package schedule

import (
	"fmt"
	"time"
)

type GinSchedule struct {
	Cron      string `json:"cron"`
	Immediate bool   `json:"immediate"`
	Disable   bool   `json:"disable"`
	Task      func()
}

// 定时获取必应的壁纸
func GetBingImg() GinSchedule {
	task := GinSchedule{
		Cron:      "*/5 * * * * ?",
		Immediate: true,
	}
	task.Task = func() {
		fmt.Println("hello world", time.Now())
	}
	return task
}
