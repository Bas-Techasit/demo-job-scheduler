package services

type JobService interface {
	NewJob(NewJob) error
	ScheduleAllJob()
}

type NewJob struct {
	JobName string
	CronExp string
}
