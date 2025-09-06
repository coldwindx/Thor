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
		ctx.Router.POST("/job/insert", Insert)
		ctx.Router.POST("/job/query", Query)
	})
}

func Insert(c *gin.Context) {
	var form request.JobInsertReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}
	var job models.Job
	_ = deepcopier.Copy(&form).To(&job)
	_, err := services.JobService.Insert(job)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, job)
}

func Query(c *gin.Context) {
	var form request.JobQueryReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}

	var query models.JobQuery
	_ = deepcopier.Copy(&form).To(&query)
	jobs, err := services.JobService.Query(query)

	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, jobs)
}
