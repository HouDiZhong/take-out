package initialize

import (
	"take-out/tasks"
)

func initCron() {
	var s = tasks.NewScheduler()

	s.AddTask("ProcessTimeoutOrder", "0/1 * * * *", tasks.ProcessTimeoutOrder)
	s.Start()
}
