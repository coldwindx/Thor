package job

func init() {
	executor := TaskInputJobExecutor{}
	JobScheduler.ExecutorManager[executor.GetName()] = &executor
}

type TaskInputJobExecutor struct {
	IJobExecutor
}

func (it *TaskInputJobExecutor) GetName() string {
	return "task_input_executor"
}
