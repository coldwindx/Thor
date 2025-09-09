package job

import (
	"Thor/utils"
	"fmt"
)

var JobSchedulerImpl = &JobScheduler{manager: make(map[string]IJobExecutor)}

func init() {
	fmt.Println("init job scheduler")
	utils.ScanInject("JobScheduler", JobSchedulerImpl)
}

type JobScheduler struct {
	manager map[string]IJobExecutor
}

func (scheduler *JobScheduler) Register(executor IJobExecutor) {
	scheduler.manager[executor.GetName()] = executor
}

func (scheduler *JobScheduler) GetExecutor(name string) IJobExecutor {
	return scheduler.manager[name]
}
