package test

import (
	"Thor/bootstrap"
	"Thor/src/mapper"
	"Thor/utils/invoke"
	proxy2 "Thor/utils/proxy"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestProxyStruct(t *testing.T) {
	cat := &Cat{Name: "kitty"}
	proxy := proxy2.NewMethodProxy(&DefaultAnimal{}, func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
		// 打印方法名和参数
		fmt.Println("[method] >>>", method.Name, args)
		return method.Invoke(cat, args)
	}).(*DefaultAnimal)
	// 调用代理方法
	name := proxy.GetName()
	assert.Equal(t, "cat:kitty", name)
}

func TestProxyMapper(t *testing.T) {
	// 指定环境变量
	_ = os.Setenv("VIPER_CONFIG", "../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()
	bootstrap.Beans.Populate()

	jobMapper := bootstrap.Beans.GetByName("JobMapper").(*mapper.JobMapper)
	assert.Equal(t, "JobMapper.Test()", jobMapper.Test())
}
