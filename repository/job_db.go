package repository

import (
	"github.com/jmoiron/sqlx"
)

type jobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) JobRepository {
	return jobRepository{db: db}
}

func (r jobRepository) GetAll() ([]Job, error) {
	jobs := []Job{}
	query := `
		SELECT job_id, job_name, cron_exp, create_date 
		FROM jobs
	`
	err := r.db.Select(&jobs, query)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r jobRepository) GetById(jobID string) (*Job, error) {
	return nil, nil
}

func (r jobRepository) Create(job Job) error {
	query := `
		INSERT INTO jobs(job_id, job_name, cron_exp, create_date) 
		VALUE (?, ?, ?, ?)
	`
	r.db.Exec(
		query,
		job.JobID,
		job.JobName,
		job.CronExp,
		job.CreateDate,
	)
	return nil
}

func (r jobRepository) Update(Job) (*Job, error) {
	return nil, nil
}

func (r jobRepository) Delete(string) error {
	return nil
}
