package job

import (
	"Thor/src/models"
)

func init() {
	SchedulerImpl.Register(&IosJobExecutor{})
}

type IosJobExecutor struct {
	IJobExecutor
}

func (it *IosJobExecutor) GetName() string {
	return "ios_job"
}

func (it *IosJobExecutor) Unlock(job *models.Job, args *map[string]string, options *map[string]string) error {
	return nil
}

func (it *IosJobExecutor) Start(job *models.Job) error {
	return nil
}

func (it *IosJobExecutor) Callback() error {
	return nil
}
