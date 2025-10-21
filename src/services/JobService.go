package services

import (
	"Thor/bootstrap"
	"Thor/src/mapper"
	"Thor/src/models"
	"Thor/src/models/do"
	"Thor/utils"
	"Thor/utils/inject"
	"Thor/utils/proxy"
	"context"
)

func init() {
	proxy.CycleProvide(bootstrap.Beans, &inject.Object{Name: "JobServiceImpl", Value: new(JobServiceImpl)})
}

type JobService struct {
	Create func(ctx context.Context, job *models.Job) (int64, error) `transaction:""`
}

type JobServiceImpl struct {
	JobService *JobService       `bean:"JobService;proxy"`
	JobMapper  *mapper.JobMapper `inject:"JobMapper"`
}

func (jsl *JobServiceImpl) Create(ctx context.Context, job *models.Job) (int64, error) {
	// 转换为DO对象
	jobDo := &do.Job{}
	if err := utils.Copy(job).To(jobDo); err != nil {
		return 0, err
	}
	// 插入job
	id, err := jsl.JobMapper.Insert(ctx, jobDo)
	return id, err
}
