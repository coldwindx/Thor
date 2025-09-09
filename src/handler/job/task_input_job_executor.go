package job

func init() {
	JobSchedulerImpl.Register(&TaskInputJobExecutor{})
}

type TaskInputJobExecutor struct {
	IJobExecutor
}

func (it *TaskInputJobExecutor) GetName() string {
	return "task_input_job"
}
