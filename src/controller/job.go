package controller

import (
	"Thor/common"
	"Thor/ctx"
	"Thor/src/services"
	"github.com/gin-gonic/gin"
)

func init() {
	ctx.Routes = append(ctx.Routes, func(r *gin.Engine) {
		ctx.Router.POST("/job/query", Query)
	})
}

func Query(c *gin.Context) {
	jobs, err := services.JobService.Query()
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, jobs)
}
