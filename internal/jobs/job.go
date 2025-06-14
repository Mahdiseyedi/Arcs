package jobs

import (
	"github.com/robfig/cron/v3"
)

type Job struct {
	C *cron.Cron
}

func NewCronJob() *Job {
	c := cron.New()
	return &Job{
		C: c,
	}
}

func (j *Job) Start() {
	j.C.Start()
}

func (j *Job) Stop() {
	ctx := j.C.Stop()
	<-ctx.Done()
}
