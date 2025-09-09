package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
)

var TaskMapperImpl = new(TaskMapper)

func init() {
	ctx.MybatisMapperBinds = append(ctx.MybatisMapperBinds, ctx.MybatisMapperBind{
		XmlFile: "/mapper/TaskMapper.xml",
		Mapper:  TaskMapperImpl,
	})
}

type TaskMapper struct {
	Insert func(task models.Task) (int, error)
}
