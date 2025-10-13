package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
	"Thor/utils/inject"
)

func init() {
	var TaskMapperImpl = new(TaskMapper)
	// 注入mapper
	inject.ScanInject("TaskMapper", TaskMapperImpl)
	ctx.MybatisMapperBinds = append(ctx.MybatisMapperBinds, ctx.MybatisMapperBind{
		XmlFile: "/mapper/TaskMapper.xml",
		Mapper:  TaskMapperImpl,
	})
}

type TaskMapper struct {
	Insert func(task models.Task) (int, error)
	Query  func(task models.TaskQuery) ([]models.Task, error)
}
