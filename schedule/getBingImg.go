package schedule

import (
	"fmt"
	"gin_app/util"
	"time"
)

func GetBingImg() util.GinSchedule {
	task := util.GinSchedule{
		Cron:      "*/5 * * * * ?",
		Immediate: true,
		// Disable: true,
	}
	task.Task = func() {
		fmt.Println("hello world", time.Now())
	}
	return task
}
