package services

import (
	"fmt"
	"gocron/logs"
	"gocron/repository"
	"gocron/work"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
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
		panic(err)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	cron := gocron.NewScheduler(loc)
	cron.SingletonModeAll()

	for _, job := range jobs {
		provideJob(job, cron)
	}

	cron.StartAsync()

	for cron.Len() > 0 {
		// *** Handle New Jobs And update Job *** //
		newJobs, _ := s.jobRepo.GetByPeriodTime(1, true)
		if len(newJobs) > 0 {
			for _, newJob := range newJobs {
				updatedJob, _ := cron.FindJobsByTag(strconv.Itoa(newJob.JobID))
				if len(updatedJob) != 0 {
					if updatedJob[0].IsRunning() {
						fmt.Println("continue")
						continue
					}
					tag, _ := strconv.Atoi(updatedJob[0].Tags()[0])
					for index, currentJob := range jobs {
						if currentJob.JobID == tag {
							jobs = append(jobs[:index], jobs[index+1:]...) // remove old job
							fmt.Println("updated job")
							provideJob(newJob, cron)
							jobs = append(jobs, newJob)
							break
						}
					}
				} else {
					fmt.Println("has a new job")
					provideJob(newJob, cron)
					jobs = append(jobs, newJob)
				}
			}
		}

		inactiveJobs, _ := s.jobRepo.GetByPeriodTime(1, false)
		if len(inactiveJobs) != 0 {
			var tag string
			for _, job := range inactiveJobs {
				tag = strconv.Itoa(job.JobID)
				j, _ := cron.FindJobsByTag(tag)
				if len(j) != 0 && !j[0].IsRunning() {
					fmt.Printf("job tag: %v removed\n", tag)
					cron.RemoveByTag(tag)
					// removed from slice jobs
					for index, currentJob := range jobs {
						if currentJob.JobID == job.JobID {
							jobs = append(jobs[:index], jobs[index+1:]...)
							fmt.Printf("remove tag: %v from slice", tag)
							break
						}
					}
				}
			}
		}
		fmt.Println("End Loop")
		fmt.Println(jobs)
		time.Sleep(30 * time.Second)

	}
}
func provideJob(job repository.Job, cron *gocron.Scheduler) {
	tag := strconv.Itoa(job.JobID)
	cron.RemoveByTag(tag)

	var task func()
	switch job.JobCode {
	case "job_1":
		task = func() {
			work.RunJob1()
		}
	case "job_2":
		task = func() {
			work.RunJob2()
		}
	case "job_3":
		task = func() {
			work.RunJob3()
		}
	default:
		logs.Info("Unknow Job")
		task = nil
	}
	if task != nil {
		// cron.Cron(job.ScheduleExp).Tag(tag).Do(task)
		cron.Every(20).Tag(tag).Do(task)
	}
}
