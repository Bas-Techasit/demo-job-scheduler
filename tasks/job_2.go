package tasks

import "gocron/logs"

func RunJob2(loc string) {
	logs.Info("job_2 is running from " + loc)
}
