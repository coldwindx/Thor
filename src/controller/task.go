package controller

import (
	"Thor/common"
	"Thor/ctx"
	"Thor/src/models"
	"Thor/src/request"
	"Thor/src/services"
	"Thor/tools"
	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"
)

func init() {
	ctx.Routes = append(ctx.Routes, func(r *gin.Engine) {
		ctx.Router.POST("/task/create", Create)
	})
}

func Create(c *gin.Context) {
	var form request.TaskCreateReq
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}

	var task models.Task
	_ = deepcopier.Copy(&form).To(&task)
	_, err := services.TaskService.Create(&task)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, task)
}
