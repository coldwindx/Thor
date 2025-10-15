package controller

import (
	"Thor/bootstrap"
	"Thor/common"
	"Thor/src/models"
	"Thor/src/request"
	"Thor/src/services"
	"Thor/tools"
	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"
)

func init() {
	var TaskControllerImpl = new(TaskController)
	// 路由注入
	bootstrap.Routes = append(bootstrap.Routes, func(r *gin.Engine) {
		bootstrap.Router.POST("/task/create", TaskControllerImpl.Create)
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
	if err := it.TaskService.Create(&task); nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, task)
}
