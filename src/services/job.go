package services

import (
	"Thor/src/mapper"
	"Thor/src/models"
)

type jobService struct {
}

var JobService = new(jobService)

func (it *jobService) Query() ([]models.Job, error) {
	//jobs, err := mapper.JobMapper.SelectTemplate("test")
	jobs, err := mapper.JobMapper.Query(models.Job{Name: ""})
	return jobs, err

}
