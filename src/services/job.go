package services

import (
	"Thor/ctx"
	"Thor/src/mapper"
	"Thor/src/models"
	"time"
)

type jobService struct {
}

var JobService = new(jobService)

func (it *jobService) Insert(job models.Job) (int, error) {
	it.beforeInsert(&job)
	return mapper.JobMapper.Insert(job)
}

func (it *jobService) Query(query models.JobQuery) ([]models.Job, error) {
	return mapper.JobMapper.Query(query)
}

func (it *jobService) beforeInsert(job *models.Job) {
	job.Id = ctx.Snowflake.Generate().Int64()
	t := time.Now()
	if job.CreatedAt.IsZero() {
		job.CreatedAt = t
	}
	if job.AwakenAt.IsZero() {
		job.AwakenAt = t
	}
	if job.UpdatedAt.IsZero() {
		job.UpdatedAt = t
	}
}
