package aop

import (
	"Thor/bootstrap"
	_ "Thor/bootstrap/aop/transaction"
	_ "Thor/bootstrap/initializer"
	"Thor/src/models"
	"Thor/src/services"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTransactionAop(t *testing.T) {
	// 指定环境变量
	_ = os.Setenv("VIPER_CONFIG", "../../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()

	jobService := bootstrap.Beans.GetByName("JobService").(*services.JobService)
	// 测试事务切面
	job := &models.Job{
		Name: "test",
	}
	id, err := jobService.Create(context.Background(), job)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)
}
