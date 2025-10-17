package mapper

import (
	"Thor/bootstrap"
	"Thor/src/models"
	"Thor/utils/inject"
	"context"
	"gorm.io/gorm"
)

func init() {
	bootstrap.Beans.CycleProvide(&inject.Object{Name: "JobMapperImpl", Value: new(JobMapperImpl)})
}

type JobMapper struct {
	Test   func(ctx context.Context) string `transaction:"read"`
	Insert func(job models.Job) (int, error)
	Query  func(ctx context.Context, query *models.JobQuery) ([]*models.Job, error)
	Delete func(jobQuery models.JobQuery) (int, error)
}

type JobMapperImpl struct {
	JobMapper *JobMapper `bean:"JobMapper;proxy"`
}

func (j *JobMapperImpl) Test(ctx context.Context) string {
	return "JobMapper.Test()"
}

func (j *JobMapperImpl) Query(ctx context.Context, query *models.JobQuery) ([]models.Job, error) {
	db := bootstrap.Beans.GetByName("MysqlConnection").(*gorm.DB)
	find, err := gorm.G[models.Job](db).Where(query).Find(ctx)
	if err != nil {
		return nil, err
	}
	return find, nil
}
