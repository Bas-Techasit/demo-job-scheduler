package tasks

import "gocron/logs"

func RunJob3(loc string) {
	logs.Info("job_3 is running from " + loc)
}
