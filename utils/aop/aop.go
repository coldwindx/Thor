package aop

import (
	"Thor/bootstrap"
	"Thor/utils/inject"
	"Thor/utils/invoke"
	"reflect"
)

func init(){
	bootstrap.Beans.Provide(&inject.Object{Name: })
}

// AspectManager 统一的切面管理类
type AspectManager struct {
}

// Aspect 切面接口
type Aspect interface {
	Around(obj any, method *invoke.Method, args []reflect.Value)
}
