package services

import (
	"Thor/bootstrap"
	"Thor/src/manager"
	"Thor/src/models"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"io"
	"strings"
	"time"
)

func init() {

}

type TaskService struct {
	TaskManager *manager.TaskManager `inject:"TaskManager"`
}

func (it *TaskService) Create(task *models.Task) error {
	it.beforeInsert(task)
	// read pipelines.json
	pipelines, err := it.parsePipeline()
	if err != nil {
		return err
	}
	pipe := lo.Filter(pipelines, func(item models.Pipeline, index int) bool { return item.Name == task.Pipeline })[0]
	// Workflow to DAG
	var jobMap = make(map[string]*models.Job)
	var dag = make(map[string]*models.WorkNode)

	for _, link := range strings.Split(pipe.Workflow, ",") {
		edges := strings.SplitN(link, "->", 2)
		for _, edge := range edges {
			if _, ok := dag[edge]; ok {
				continue
			}

			j := new(models.Job)
			j.Id = bootstrap.Snowflake.Generate().Int64()
			j.TaskId = task.Id
			j.Name = edge
			j.Status = int(models.INITING)
			j.AwakenAt = time.Now()
			j.CreatedAt = time.Now()
			j.UpdatedAt = time.Now()

			jobMap[j.Name] = j
			dag[j.Name] = &models.WorkNode{JobId: j.Id, JobName: j.Name}

		}

		dag[edges[0]].Next = append(dag[edges[0]].Next, dag[edges[1]].JobId)
		jobMap[edges[1]].Locked += 1
		jobMap[edges[1]].Status = int(models.WAITING)
	}

	workNodes := lo.MapToSlice(dag, func(_ string, value *models.WorkNode) *models.WorkNode { return value })
	if task.Dag, err = jsoniter.MarshalToString(workNodes); err != nil {
		return err
	}
	jobs := lo.MapToSlice(jobMap, func(_ string, value *models.Job) *models.Job { return value })
	return it.TaskManager.Create(task, jobs)
}

func (it *TaskService) beforeInsert(task *models.Task) {
	task.Id = bootstrap.Snowflake.Generate().Int64()
	t := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = t
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = t
	}
}

func (it *TaskService) parsePipeline() ([]models.Pipeline, error) {
	file, err := bootstrap.Statik.Open("/pipelines.json")
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
