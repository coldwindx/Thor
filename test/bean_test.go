package test

import (
	"Thor/bootstrap"
	"Thor/ctx"
	"Thor/src/services"
	"Thor/utils/inject"
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
	_ = g.Populate()

	zoo := g.GetByName("zoo").(*Zoo)
	assert.Equal(t, "cat:", zoo.Cat.GetName())
}

func TestServiceProxyWithInject(t *testing.T) {
	// 指定环境变量
	_ = os.Setenv("VIPER_CONFIG", "../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()

	jobServiceImpl := ctx.Beans.GetByName("JobServiceImpl").(*services.JobServiceImpl)
	assert.Equal(t, "JobServiceImpl->JobMapper.Test()", jobServiceImpl.Test())

	jobService := ctx.Beans.GetByName("JobService").(*services.JobService)
	assert.Equal(t, "JobServiceImpl->JobMapper.Test()", jobService.Test())
}
