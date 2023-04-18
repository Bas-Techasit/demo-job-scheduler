package services

import (
	"errors"
	"fmt"
	"gocron/logs"
	"gocron/repository"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron"
)

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) JobService {
	return jobService{jobRepo: jobRepo}
}

func (s jobService) ScheduleAllJob() {
	jobs, err := s.jobRepo.GetAll()
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	c := cron.New()
	for _, job := range jobs {
		c.AddFunc(job.CronExp, func() {
			logs.Info(fmt.Sprintf("JobID: %s, JobName: %s", job.JobID, job.JobName))
		})
	}
	c.Start()
	time.Sleep(2 * time.Hour)
	c.Stop()
}

func (s jobService) NewJob(newJob NewJob) error {
	if newJob.JobName == "" || newJob.CronExp == "" {
		return errors.New("job name or cron expression is empty")
	}

	job := repository.Job{
		JobID:      uuid.NewString(),
		JobName:    newJob.JobName,
		CronExp:    newJob.CronExp,
		CreateDate: time.Now(),
	}

	err := s.jobRepo.Create(job)
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	return nil
}
