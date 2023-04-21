package services

import (
	"github.com/go-co-op/gocron"
)

type JobService interface {
	ScheduleAllJob(*gocron.Scheduler, string)
}

type NewJob struct {
	JobName string
	CronExp string
}
