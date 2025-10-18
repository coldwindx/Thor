package service

import (
	"Thor/bootstrap"
	"Thor/src/models"
	"Thor/src/services"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestJobInsert(t *testing.T) {
	_ = os.Setenv("VIPER_CONFIG", "../../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()

	jobService := bootstrap.Beans.GetByName("JobService").(*services.JobService)
	job := &models.Job{Name: "test"}
	id, err := jobService.Create(context.Background(), job)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)
}
