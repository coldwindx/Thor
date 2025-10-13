package job

import (
	"Thor/src/mapper"
	"Thor/src/models"
)

func init() {
	SchedulerImpl.Register(&IosJobExecutor{})
}

type IosJobExecutor struct {
	IJobExecutor
	JobMapper *mapper.JobMapper `inject:""`
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

func (it *IosJobExecutor) Callback(body string) (any, error) {
	return nil, nil
}
