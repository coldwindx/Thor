package job

import (
	"Thor/src/models"
)

func init() {
	SchedulerImpl.Register(&TaskInputJobExecutor{})
}

type TaskInputJobExecutor struct {
	IJobExecutor
}

func (it *TaskInputJobExecutor) GetName() string {
	return "task_input_job"
}

func (it *TaskInputJobExecutor) Unlock(job *models.Job, args *map[string]string, options *map[string]string) error {
	return nil
}

func (it *TaskInputJobExecutor) Start(job *models.Job) error {
	return nil
}

func (it *TaskInputJobExecutor) Callback() error {
	return nil
}
