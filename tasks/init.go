package tasks

import (
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	*gocron.Scheduler
	taskIds map[string]*gocron.Job
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Scheduler: gocron.NewScheduler(time.Local),
		taskIds:   make(map[string]*gocron.Job),
	}
}

func (s *Scheduler) AddTask(name, cron string, jobfun func()) {
	if _, ok := s.taskIds[name]; !ok {
		job, _ := s.Cron(cron).Do(jobfun)
		s.taskIds[name] = job
	}
}

func (s *Scheduler) RemoveTask(name string) {
	if job, ok := s.taskIds[name]; ok {
		s.RemoveByReference(job)
		delete(s.taskIds, name)
	}
}

func (s *Scheduler) Start() {
	s.StartAsync()
}
