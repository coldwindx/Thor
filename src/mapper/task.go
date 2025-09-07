package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
)

func init() {
	ctx.MybatisMapperBinds = append(ctx.MybatisMapperBinds, ctx.MybatisMapperBind{
		XmlFile: "/mapper/TaskMapper.xml",
		Mapper:  TaskMapper,
	})
}

var TaskMapper = new(taskMapper)

type taskMapper struct {
	Insert func(task models.Task) (int, error)
}
