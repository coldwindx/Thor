package mapper

import (
	"Thor/bootstrap"
	"Thor/src/models"
	"Thor/utils/inject"
	"context"
	"gorm.io/gorm"
)

var JobMapperImpl = new(JobMapper)

func init() {
	bootstrap.Beans.Provide(&inject.Object{Name: "JobMapper", Value: JobMapperImpl})
}

type JobMapper struct {
	Insert func(job models.Job) (int, error)
	Delete func(jobQuery models.JobQuery) (int, error)
}

func (j *JobMapper) Test() string {
	return "JobMapper.Test()"
}

func (j *JobMapper) Query(ctx context.Context, query *models.JobQuery) ([]models.Job, error) {
	db := bootstrap.Beans.GetByName("MysqlConnection").(*gorm.DB)
	find, err := gorm.G[models.Job](db).Where(query).Find(ctx)
	if err != nil {
		return nil, err
	}
	return find, nil
}
