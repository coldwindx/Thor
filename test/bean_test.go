package test

import (
	"Thor/bootstrap"
	"Thor/src/services"
	"Thor/utils/inject"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBeanAutoProvide(t *testing.T) {
	type Zoo struct {
		Cat *Cat `bean:"cat"`
	}

	g := inject.NewGraph()
	g.CycleProvide(&inject.Object{Value: &Zoo{}})
	assert.Equal(t, "cat:", g.GetByName("cat").(*Cat).GetName())
}

func TestBeanAutoPopulate(t *testing.T) {
	type Zoo struct {
		Cat *Cat `bean:"cat" inject:"cat"`
	}

	g := inject.NewGraph()
	g.CycleProvide(&inject.Object{Name: "zoo", Value: &Zoo{}})
	g.Populate()

	zoo := g.GetByName("zoo").(*Zoo)
	assert.Equal(t, "cat:", zoo.Cat.GetName())
}

func TestServiceProxyWithInject(t *testing.T) {
	// 指定环境变量
	_ = os.Setenv("VIPER_CONFIG", "../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()

	jobServiceImpl := bootstrap.Beans.GetByName("JobServiceImpl").(*services.JobServiceImpl)
	assert.Equal(t, "JobServiceImpl->JobMapper.Test()", jobServiceImpl.Test(context.Background()))

	jobService := bootstrap.Beans.GetByName("JobService").(*services.JobService)
	assert.Equal(t, "JobServiceImpl->JobMapper.Test()", jobService.Test())
}
