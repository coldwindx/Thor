package services

import (
	"Thor/ctx"
	"Thor/src/handler/job"
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/utils"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/zhuxiujia/GoMybatis"
	"io"
	"strings"
	"time"
)

func init() {
	var impl = new(TaskServiceImpl)
	impl.JobScheduler = job.SchedulerImpl
	impl.TaskService.Create = impl.Create
	// bean注入
	utils.ScanInject("TaskServiceImpl", impl)
	// aop必须注在接口上
	GoMybatis.AopProxyService(&impl.TaskService, &ctx.MybatisEngine)
}

type TaskService struct {
	Create func(task *models.Task) (int, error) `tx:"" rollback:"error"`
}

type TaskServiceImpl struct {
	TaskService  `bean:"TaskService"`
	JobScheduler *job.Scheduler     `inject:"JobScheduler"`
	TaskMapper   *mapper.TaskMapper `inject:"TaskMapper"`
	JobMapper    *mapper.JobMapper  `inject:"JobMapper"`
}

func (it *TaskServiceImpl) Create(task *models.Task) (int, error) {
	it.beforeInsert(task)
	// read pipelines.json
	pipelines, err := it.parsePipeline()
	if err != nil {
		return 0, err
	}
	pipe := lo.Filter(pipelines, func(item models.Pipeline, index int) bool { return item.Name == task.Pipeline })[0]
	// Workflow to DAG
	var jobs = make([]*models.Job, 0)
	var dag = make(map[string]*models.WorkNode)

	for _, link := range strings.Split(pipe.Workflow, ",") {
		edges := strings.SplitN(link, "->", 2)
		for _, edge := range edges {
			if _, ok := dag[edge]; ok {
				continue
			}

			j := new(models.Job)
			j.Id = ctx.Snowflake.Generate().Int64()
			j.TaskId = task.Id
			j.Name = edge
			j.AwakenAt = time.Now()
			j.CreatedAt = time.Now()
			j.UpdatedAt = time.Now()
			jobs = append(jobs, j)
			dag[j.Name] = &models.WorkNode{JobId: j.Id, JobName: j.Name}

			executor := it.JobScheduler.GetExecutor(j.Name)
			if executor == nil {
				return 0, errors.New("Job executor not found, " + j.Name)
			}
		}

		dag[edges[0]].Next = append(dag[edges[0]].Next, dag[edges[1]].JobId)
	}

	workNodes := lo.MapToSlice(dag, func(_ string, value *models.WorkNode) *models.WorkNode { return value })
	if task.Dag, err = jsoniter.MarshalToString(workNodes); err != nil {
		return 0, err
	}

	// 数据库动作，后续想办法把事务控制在这个方法里
	it.create(task, jobs)
	return 1, nil
}

func (it *TaskServiceImpl) create(task *models.Task, jobs []*models.Job) {
	lo.ForEach(jobs, func(job *models.Job, index int) {
		if _, err := it.JobMapper.Insert(*job); err != nil {
			panic("Job insert error, " + err.Error())
		}
	})
	if _, err := it.TaskMapper.Insert(*task); err != nil {
		panic("Task insert error, " + err.Error())
	}
}

func (it *TaskServiceImpl) beforeInsert(task *models.Task) {
	task.Id = ctx.Snowflake.Generate().Int64()
	t := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = t
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = t
	}
}

func (it *TaskServiceImpl) parsePipeline() ([]models.Pipeline, error) {
	file, err := ctx.Statik.Open("/pipelines.json")
	if err != nil {
		return nil, err
	}
	all, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var pipelines []models.Pipeline
	err = jsoniter.Unmarshal(all, &pipelines)
	if err != nil {
		return nil, err
	}
	return pipelines, err
}
