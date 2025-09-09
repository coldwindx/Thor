package job

import (
	"Thor/src/mapper"
	"Thor/utils"
	"fmt"
)

var SchedulerImpl = &Scheduler{
	manager: make(map[string]IJobExecutor),
	//JobMapper: mapper.JobMapperImpl,
}

func init() {
	fmt.Println("init job scheduler")
	utils.ScanInject("JobScheduler", SchedulerImpl)
}

type Scheduler struct {
	manager   map[string]IJobExecutor
	JobMapper *mapper.JobMapper `inject:"JobMapper"`
}

func (scheduler *Scheduler) Register(executor IJobExecutor) {
	scheduler.manager[executor.GetName()] = executor
}

func (scheduler *Scheduler) GetExecutor(name string) IJobExecutor {
	return scheduler.manager[name]
}

func (scheduler *Scheduler) Unlock(id int64) error {
	return nil
}

func (scheduler *Scheduler) Start(id int64) error {
	return nil
}

func (scheduler *Scheduler) Callback(id int64) error {
	return nil
}
