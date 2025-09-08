package services

import (
	"Thor/ctx"
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/utils"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"io"
	"strings"
	"time"
)

func init() {
	fmt.Println("init bean: TaskServiceImpl...")
	var TaskServiceImpl = new(TaskService)
	utils.ScanInject("TaskServiceImpl", TaskServiceImpl)
	//GoMybatis.AopProxyService(TaskServiceImpl, &ctx.MybatisEngine)
}

type TaskService struct {
	TaskMapper *mapper.TaskMapper `inject:"TaskMapperImpl"`
}

func (it *TaskService) Create(task *models.Task) (int, error) {
	it.beforeInsert(task)
	// read pipelines.json
	pipelines, err := it.parsePipeline()
	if err != nil {
		return 0, err
	}
	pipe := lo.Filter(pipelines, func(item models.Pipeline, index int) bool { return item.Name == task.Pipeline })[0]
	// Workflow to DAG
	var jobs = make(map[string]*models.Job)
	var dag = make(map[string]*models.WorkNode)

	for _, link := range strings.Split(pipe.Workflow, ",") {
		edges := strings.SplitN(link, "->", 2)
		for _, edge := range edges {
			if _, ok := jobs[edge]; ok {
				continue
			}

			job := new(models.Job)
			job.Id = ctx.Snowflake.Generate().Int64()
			job.TaskId = task.Id
			job.Name = edge
			job.AwakenAt = time.Now()
			job.CreatedAt = time.Now()
			job.UpdatedAt = time.Now()

			dag[job.Name] = &models.WorkNode{JobId: job.Id, JobName: job.Name}
			jobs[job.Name] = job
		}

		dag[edges[0]].Next = append(dag[edges[0]].Next, dag[edges[1]].JobId)
	}

	workNodes := lo.MapToSlice(dag, func(_ string, value *models.WorkNode) *models.WorkNode { return value })
	if task.Dag, err = jsoniter.MarshalToString(workNodes); err != nil {
		return 0, err
	}

	// 事务
	tx, err := ctx.DefaultSqlDB.Begin()
	if err != nil {
		return 0, err
	}
	_, _ = it.TaskMapper.Insert(*task)
	err = tx.Commit()
	return 0, err
}

func (it *TaskService) beforeInsert(task *models.Task) {
	task.Id = ctx.Snowflake.Generate().Int64()
	t := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = t
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = t
	}
}

func (it *TaskService) parsePipeline() ([]models.Pipeline, error) {
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
