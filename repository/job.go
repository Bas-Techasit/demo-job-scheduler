package repository

import "time"

type Job struct {
	JobID      string    `db:"job_id"`
	JobName    string    `db:"job_name"`
	CronExp    string    `db:"cron_exp"`
	CreateDate time.Time `db:"create_date"`
}

type JobRepository interface {
	GetAll() ([]Job, error)
	GetById(string) (*Job, error)
	Create(Job) error
	Update(Job) (*Job, error)
	Delete(string) error
}
