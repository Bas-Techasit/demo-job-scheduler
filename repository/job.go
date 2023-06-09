package repository

type Job struct {
	JobID       int    `db:"job_id"`
	JobCode     string `db:"job_code"`
	ScheduleExp string `db:"schedule_exp"`
	Location    string `db:"location"`
}

type JobRepository interface {
	GetAll() ([]Job, error)
	GetByPeriodTime(float64, bool) ([]Job, error)
}
