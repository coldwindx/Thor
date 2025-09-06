package services

import (
	"Thor/src/mapper"
	"Thor/src/models"
)

type jobService struct {
}

var JobService = new(jobService)

func (it *jobService) Insert(job models.Job) (int, error) {
	return mapper.JobMapper.Insert(job)
}

func (it *jobService) Query(query models.JobQuery) ([]models.Job, error) {
	return mapper.JobMapper.Query(query)
}
