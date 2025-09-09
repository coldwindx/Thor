package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
)

var JobMapperImpl = new(JobMapper)

func init() {
	ctx.MybatisMapperBinds = append(ctx.MybatisMapperBinds, ctx.MybatisMapperBind{
		XmlFile: "/mapper/JobMapper.xml",
		Mapper:  JobMapperImpl,
	})
}

type JobMapper struct {
	Insert func(job models.Job) (int, error)
	Delete func(jobQuery models.JobQuery) (int, error)
	Query  func(jobQuery models.JobQuery) ([]models.Job, error)
}
