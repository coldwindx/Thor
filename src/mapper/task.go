package mapper

import (
	"Thor/ctx"
	"Thor/src/models"
	"Thor/utils"
	"fmt"
)

func init() {
	fmt.Println("init bean: TaskMapperImpl...")
	var TaskMapperImpl = new(TaskMapper)
	utils.ScanInject("TaskMapperImpl", TaskMapperImpl)
	ctx.MybatisMapperBinds = append(ctx.MybatisMapperBinds, ctx.MybatisMapperBind{
		XmlFile: "/mapper/TaskMapper.xml",
		Mapper:  TaskMapperImpl,
	})
}

type TaskMapper struct {
	Insert func(task models.Task) (int, error)
}
