package schedule

import (
	"fmt"
	"net/http"
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
		Cron:      "0 5 0 * * ?",
		Immediate: true,
	}
	task.Task = func() {
		_, err := http.Get("http://127.0.0.1:8000/api/img?schedule=1")
		if err != nil {
			fmt.Println("img err:", err)
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " 必应壁纸下载成功")
	}
	return task
}
