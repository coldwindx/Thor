package mapper

import (
	"Thor/bootstrap"
	"Thor/src/models/do"
	"Thor/utils/inject"
	"context"
	"gorm.io/gorm"
)

func init() {
	bootstrap.Beans.CycleProvide(&inject.Object{Name: "JobMapperImpl", Value: new(JobMapperImpl)})
}

type JobMapper struct {
	Test   func(ctx context.Context) string `transaction:"read"`
	Insert func(ctx context.Context, job *do.Job) (int64, error)
	Query  func(ctx context.Context, query *do.JobQuery) ([]*do.Job, error)
	Delete func(jobQuery do.JobQuery) (int, error)
}

type JobMapperImpl struct {
	Db        *bootstrap.DBClient `inject:"DBClient"`
	JobMapper *JobMapper          `bean:"JobMapper;proxy"`
}

func (j *JobMapperImpl) Test(ctx context.Context) string {
	return "JobMapper.Test()"
}

func (j *JobMapperImpl) Insert(ctx context.Context, job *do.Job) (int64, error) {
	// 插入job
	err := j.Db.DB(ctx).Create(job).Error
	return job.ID, err
}

func (j *JobMapperImpl) Query(ctx context.Context, query *do.JobQuery) ([]*do.Job, error) {
	db := ctx.Value("sql").(*gorm.DB)
	find, err := gorm.G[*do.Job](db).Where(query).Find(ctx)
	if err != nil {
		return nil, err
	}
	return find, nil
}
