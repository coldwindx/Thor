package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
)

func init() {
	ctx.MybatisMapperBinds = append(ctx.MybatisMapperBinds, ctx.MybatisMapperBind{
		XmlFile: "/mapper/JobMapper.xml",
		Mapper:  JobMapper,
	})
}

var JobMapper = new(jobMapper)

type jobMapper struct {
	Insert func(job models.Job) (int, error)
	Query  func(jobQuery models.JobQuery) ([]models.Job, error)
}
