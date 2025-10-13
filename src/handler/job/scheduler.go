package job

import (
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/utils"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
)

var SchedulerImpl = &Scheduler{manager: make(map[string]IJobExecutor)}

func init() {
	fmt.Println("init job scheduler")
	utils.ScanInject("JobScheduler", SchedulerImpl)
}

type Scheduler struct {
	manager    map[string]IJobExecutor
	TaskMapper mapper.TaskMapper `inject:""`
	JobMapper  mapper.JobMapper  `inject:""`
}

func (it *Scheduler) Register(executor IJobExecutor) {
	it.manager[executor.GetName()] = executor
}

func (it *Scheduler) GetExecutor(name string) IJobExecutor {
	return it.manager[name]
}

func (it *Scheduler) Unlock(id int64) error {
	return nil
}

func (it *Scheduler) Start(job *models.Job) error {
	executor, ok := it.manager[job.Name]
	if !ok {
		return errors.New("not found job executor, name is : " + job.Name)
	}
	return executor.Start(job)
}

func (it *Scheduler) Callback(body string) error {
	var name string
	// step1 存储output
	executor, ok := it.manager[name]
	if !ok {
		return errors.New("not found job executor, name is : " + name)
	}
	_, err := executor.Callback(body)
	if err != nil {
		return err
	}
	// step2 寻找后继jobs
	query := models.JobQuery{PageSize: 1}
	jobs, err := it.JobMapper.Query(query)
	if err != nil {
		return err
	}
	j, ok := lo.First(jobs)
	if !ok {
		return errors.New("not found job")
	}

	taskQuery := models.TaskQuery{Id: j.TaskId, PageSize: 1}
	tasks, err := it.TaskMapper.Query(taskQuery)
	if err != nil {
		return err
	}
	task, ok := lo.First(tasks)
	if !ok {
		return errors.New("not found task")
	}
	var workflow []models.WorkNode
	if err = jsoniter.UnmarshalFromString(task.Dag, workflow); err != nil {
		return err
	}
	node, ok := lo.Find(workflow, func(item models.WorkNode) bool { return item.JobId == j.Id })
	if !ok {
		return errors.New("not found work node")
	}
	// step3 unlock后继jobs，通过消息实现
	return nil
}
