package test

import (
	"Thor/src/models"
	"Thor/src/models/do"
	"Thor/utils/beans"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBeanCopy(t *testing.T) {
	job := &models.Job{Name: "test"}
	jobDo := &do.Job{}
	err := beans.Copy(job).To(jobDo)
	assert.Nil(t, err)
	assert.Equal(t, job.Name, jobDo.Name)
}
