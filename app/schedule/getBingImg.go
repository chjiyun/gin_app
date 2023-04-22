package schedule

import (
	"fmt"
	"net/http"
	"time"
)

// GetBingImg 定时获取必应的壁纸
func (s Schedule) GetBingImg() MySchedule {
	task := MySchedule{
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
		fmt.Println(time.Now().Format(time.DateTime), " 必应壁纸下载成功")
	}
	return task
}
