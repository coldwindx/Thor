package job

import (
	"Thor/src/models"
)

func init() {
	SchedulerImpl.Register(&VerifyJobExecutor{})
}

type VerifyJobExecutor struct {
	IJobExecutor
}

func (it *VerifyJobExecutor) GetName() string {
	return "verify_job"
}

func (it *VerifyJobExecutor) Unlock(job *models.Job, args *map[string]string, options *map[string]string) error {
	return nil
}

func (it *VerifyJobExecutor) Start(job *models.Job) error {
	return nil
}

func (it *VerifyJobExecutor) Callback() error {
	return nil
}
