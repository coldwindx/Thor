package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
	"Thor/utils/inject"
)

var JobMapperImpl = new(JobMapper)

func init() {

	ctx.Beans.Provide(&inject.Object{Name: "JobMapper", Value: JobMapperImpl})
}

type JobMapper struct {
	Insert func(job models.Job) (int, error)
	Delete func(jobQuery models.JobQuery) (int, error)
	Query  func(jobQuery models.JobQuery) ([]models.Job, error)
}

func (j *JobMapper) Test() string {
	return "JobMapper.Test()"
}
