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
	SelectTemplate func(name string) ([]models.Job, error) `args:"name"`
	Query          func(job models.Job) ([]models.Job, error)
}
