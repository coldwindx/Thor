package manager

import (
	"Thor/ctx"
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/utils"
	"github.com/samber/lo"
	"github.com/zhuxiujia/GoMybatis"
)

func init() {
	var impl = new(TaskManagerImpl)
	impl.TaskManager.Create = impl.Create
	utils.ScanInject("TaskManagerImpl", impl)
	GoMybatis.AopProxyService(&impl.TaskManager, &ctx.MybatisEngine)
}

type TaskManager struct {
	Create func(task *models.Task, jobs []*models.Job) `tx:"" rollback:"error"`
}

type TaskManagerImpl struct {
	TaskManager `bean:"TaskManager"`
	TaskMapper  *mapper.TaskMapper `inject:"TaskMapper"`
	JobMapper   *mapper.JobMapper  `inject:"JobMapper"`
}

func (it *TaskManagerImpl) Create(task *models.Task, jobs []*models.Job) {
	lo.ForEach(jobs, func(job *models.Job, index int) {
		if _, err := it.JobMapper.Insert(*job); err != nil {
			panic("Job insert error, " + err.Error())
		}
	})
	if _, err := it.TaskMapper.Insert(*task); err != nil {
		panic("Task insert error, " + err.Error())
	}
}
