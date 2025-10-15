package services

import (
	"Thor/bootstrap"
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/utils/inject"
	"context"
	"time"
)

func init() {
	bootstrap.Beans.CycleProvide(&inject.Object{Name: "JobServiceImpl", Value: &JobServiceImpl{}})
}

type JobService struct {
	Test   func() string
	Insert func(job *models.Job) (int, error)
	Delete func(query models.JobQuery) (int, error)
	Query  func(query *models.JobQuery) ([]models.Job, error)
}

type JobServiceImpl struct {
	JobService *JobService       `bean:"JobService;proxy"`
	JobMapper  *mapper.JobMapper `inject:"JobMapper"`
}

func (j *JobServiceImpl) Test() string {
	return "JobServiceImpl->" + j.JobMapper.Test()
}

func (it *JobServiceImpl) Insert(job *models.Job) (int, error) {
	it.beforeInsert(job)
	return it.JobMapper.Insert(*job)
}

func (it *JobServiceImpl) Delete(query models.JobQuery) (int, error) {
	return it.JobMapper.Delete(query)
}

func (it *JobServiceImpl) Query(query *models.JobQuery) ([]*models.Job, error) {
	it.beforeQuery(query)
	return it.JobMapper.Query(context.Background(), query)
}

func (it *JobServiceImpl) beforeInsert(job *models.Job) {
	job.Id = bootstrap.Snowflake.Generate().Int64()
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

func (it *JobServiceImpl) beforeQuery(query *models.JobQuery) {
	if query.CreatedAfter.IsZero() {
		query.CreatedAfter = time.Now()
	}
}
