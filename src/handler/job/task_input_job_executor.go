package job

import "fmt"

func init() {
	executor := TaskInputJobExecutor{}
	JobScheduler.ExecutorManager[executor.GetName()] = &executor
	fmt.Println("init task_input_job")
}

type TaskInputJobExecutor struct {
	IJobExecutor
}

func (it *TaskInputJobExecutor) GetName() string {
	return "task_input_job"
}
