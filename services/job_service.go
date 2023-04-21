package services

import (
	"fmt"
	"gocron/logs"
	"gocron/repository"
	"gocron/tasks"
	"strconv"

	"github.com/go-co-op/gocron"
)

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) JobService {
	return jobService{jobRepo: jobRepo}
}

func (s jobService) ScheduleAllJob(cron *gocron.Scheduler, checkTag string) {
	jobs, err := s.jobRepo.GetAll()
	if err != nil {
		panic(err)
	}
	for _, job := range jobs {
		provideJob(job, cron)
	}
	task := func() {
		checkJobsFromDB(cron, s.jobRepo, &jobs)
	}
	cron.Every(1).Minutes().Tag(checkTag).Do(task)
}

func provideJob(job repository.Job, cron *gocron.Scheduler) {
	tag := strconv.Itoa(job.JobID)
	cron.RemoveByTag(tag)

	var task func()
	switch job.JobCode {
	case "job_1":
		task = func() {
			tasks.RunJob1(job.Location)
		}
	case "job_2":
		task = func() {
			tasks.RunJob2(job.Location)
		}
	case "job_3":
		task = func() {
			tasks.RunJob3(job.Location)
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

func checkJobsFromDB(cron *gocron.Scheduler, s repository.JobRepository, jobs *[]repository.Job) {
	j := *jobs

	fmt.Println("Started check jobs ")

	newJobs, _ := s.GetByPeriodTime(1, true)
	if len(newJobs) > 0 {
		for _, newJob := range newJobs {
			updatedJob, _ := cron.FindJobsByTag(strconv.Itoa(newJob.JobID))
			if len(updatedJob) != 0 {
				if updatedJob[0].IsRunning() {
					fmt.Println("continue")
					continue
				}
				tag, _ := strconv.Atoi(updatedJob[0].Tags()[0])
				for index, currentJob := range j {
					if currentJob.JobID == tag {
						j = append(j[:index], j[index+1:]...) // remove old job
						fmt.Println("updated job")
						provideJob(newJob, cron)
						j = append(j, newJob)
						break
					}
				}
			} else {
				fmt.Println("has a new job")
				provideJob(newJob, cron)
				j = append(j, newJob)
			}
		}
	}

	inactiveJobs, _ := s.GetByPeriodTime(1, false)
	if len(inactiveJobs) != 0 {
		var tag string
		for _, job := range inactiveJobs {
			tag = strconv.Itoa(job.JobID)
			cj, _ := cron.FindJobsByTag(tag)
			if len(cj) != 0 && !cj[0].IsRunning() {
				fmt.Printf("job tag: %v removed\n", tag)
				cron.RemoveByTag(tag)
				// removed from slice jobs
				for index, currentJob := range j {
					if currentJob.JobID == job.JobID {
						j = append(j[:index], j[index+1:]...)
						fmt.Printf("remove tag: %v from slice", tag)
						break
					}
				}
			}
		}
	}

	fmt.Println("End Check jobs")
	fmt.Println(jobs)
}
