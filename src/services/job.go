package services

import (
    "Thor/ctx"
    "Thor/src/mapper"
    "Thor/src/models"
    "fmt"
    "time"
)

func init() {

}

type JobService interface {
    Test() error
}

type JobServiceImpl struct {
}

func (it *JobServiceImpl) Test() error {
    fmt.Println("Testing job service")
    return nil
}
func (it *JobServiceImpl) Insert(job *models.Job) (int, error) {
    it.beforeInsert(job)
    return mapper.JobMapperImpl.Insert(*job)
}

func (it *JobServiceImpl) Delete(query models.JobQuery) (int, error) {
    return mapper.JobMapperImpl.Delete(query)
}

func (it *JobServiceImpl) Query(query *models.JobQuery) ([]models.Job, error) {
    it.beforeQuery(query)
    return mapper.JobMapperImpl.Query(*query)
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

func (it *JobServiceImpl) beforeQuery(query *models.JobQuery) {
    if query.CreatedAfter.IsZero() {
        query.CreatedAfter = time.Now()
    }
}
