package services

import (
	"Thor/ctx"
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/utils"
	"time"
)

func init() {
	var impl = new(JobServiceImpl)
	impl.JobService.Insert = impl.Insert
	impl.JobService.Delete = impl.Delete
	impl.JobService.Query = impl.Query
	utils.ScanInject("JobServiceImpl", impl)
}

type JobService struct {
	Insert func(job *models.Job) (int, error)
	Delete func(query models.JobQuery) (int, error)
	Query  func(query models.JobQuery) ([]models.Job, error)
}

type JobServiceImpl struct {
	JobService `bean:"JobService"`
}

func (it *JobServiceImpl) Insert(job *models.Job) (int, error) {
	it.beforeInsert(job)
	return mapper.JobMapperImpl.Insert(*job)
}

func (it *JobServiceImpl) Delete(query models.JobQuery) (int, error) {
	return mapper.JobMapperImpl.Delete(query)
}

func (it *JobServiceImpl) Query(query models.JobQuery) ([]models.Job, error) {
	return mapper.JobMapperImpl.Query(query)
}

func (it *JobServiceImpl) beforeInsert(job *models.Job) {
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
