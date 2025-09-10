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
	var impl = new(JobController)
	utils.ScanInject("JobController", impl)
	ctx.Routes = append(ctx.Routes, func(r *gin.Engine) {
		ctx.Router.POST("/job/insert", impl.Insert)
		ctx.Router.POST("/job/delete", impl.Delete)
		ctx.Router.POST("/job/query", impl.Query)
	})
}

type JobController struct {
	JobService *services.JobService `inject:"JobService"`
}

func (it *JobController) Insert(c *gin.Context) {
	var form request.JobInsertReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}
	var job models.Job
	_ = deepcopier.Copy(&form).To(&job)
	_, err := it.JobService.Insert(&job)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, job)
}

func (it *JobController) Delete(c *gin.Context) {
	var form request.JobDeleteReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}

	var query models.JobQuery
	_ = deepcopier.Copy(&form).To(&query)
	jobs, err := it.JobService.Delete(query)

	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, jobs)
}

func (it *JobController) Query(c *gin.Context) {
	var form request.JobQueryReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}

	var query models.JobQuery
	_ = deepcopier.Copy(&form).To(&query)
	jobs, err := it.JobService.Query(&query)

	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, jobs)
}
