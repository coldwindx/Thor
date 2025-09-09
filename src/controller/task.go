package controller

import (
	"Thor/common"
	"Thor/ctx"
	"Thor/src/models"
	"Thor/src/request"
	"Thor/src/services"
	"Thor/tools"
	"Thor/utils"
	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"
)

func init() {
	var TaskControllerImpl = new(TaskController)
	// bean注入
	utils.ScanInject("TaskControllerImpl", TaskControllerImpl)
	// 路由注入
	ctx.Routes = append(ctx.Routes, func(r *gin.Engine) {
		ctx.Router.POST("/task/create", TaskControllerImpl.Create)
	})
}

type TaskController struct {
	TaskService *services.TaskService `inject:"TaskService"`
}

func (it *TaskController) Create(c *gin.Context) {
	var form request.TaskCreateReq
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}

	var task models.Task
	_ = deepcopier.Copy(&form).To(&task)
	_, err := it.TaskService.Create(&task)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, task)
}
