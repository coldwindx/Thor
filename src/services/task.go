package services

import (
	"Thor/ctx"
	"Thor/src/mapper"
	"Thor/src/models"
	"time"
)

var TaskService = new(taskService)

type taskService struct {
}

func (it *taskService) Create(task *models.Task) (int, error) {
	it.beforeInsert(task)
	return mapper.TaskMapper.Insert(*task)
}

func (it *taskService) beforeInsert(task *models.Task) {
	task.Id = ctx.Snowflake.Generate().Int64()
	t := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = t
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = t
	}
}
