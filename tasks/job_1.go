package tasks

import "gocron/logs"

func RunJob1(loc string) {
	logs.Info("job_1 is running from " + loc)
}
