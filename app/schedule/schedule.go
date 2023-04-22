package schedule

type Schedule struct {
}

// MySchedule 定时任务返回结构
type MySchedule struct {
	Cron      string `json:"cron"`
	Immediate bool   `json:"immediate"`
	Disable   bool   `json:"disable"`
	Task      func()
}
