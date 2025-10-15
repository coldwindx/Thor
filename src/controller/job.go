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
	var impl = new(JobController)
	bootstrap.Routes = append(bootstrap.Routes, func(r *gin.Engine) {
		bootstrap.Router.POST("/job/insert", impl.Insert)
		bootstrap.Router.POST("/job/delete", impl.Delete)
		bootstrap.Router.POST("/job/query", impl.Query)
		bootstrap.Router.POST("/job/callback", impl.Callback)
	})
}

type JobController struct {
	JobService *services.JobServiceImpl `inject:"JobService"`
}

func (it *JobController) Insert(c *gin.Context) {
	var form request.JobInsertReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}
	var j models.Job
	_ = deepcopier.Copy(&form).To(&j)
	_, err := it.JobService.Insert(&j)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, j)
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

func (it *JobController) Callback(c *gin.Context) {
	// todo: 调用callback

}
