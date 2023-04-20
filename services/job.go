package services

type JobService interface {
	ScheduleAllJob()
}

type NewJob struct {
	JobName string
	CronExp string
}
