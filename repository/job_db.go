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
		SELECT job_id, job_code, schedule_exp, location
		FROM jobs
		WHERE status = true
	`
	err := r.db.Select(&jobs, query)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r jobRepository) GetByPeriodTime(minutes float64, status bool) ([]Job, error) {
	jobs := []Job{}
	query := `
		SELECT job_id, job_code, schedule_exp
		FROM jobs
		WHERE modify_date BETWEEN SUBDATE(NOW(), INTERVAL ? MINUTE) AND NOW() AND status = ?;
	`
	err := r.db.Select(&jobs, query, minutes, status)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
