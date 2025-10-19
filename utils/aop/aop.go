package aop

import (
	"Thor/utils/invoke"
	"reflect"
)

// AspectManager 统一的切面管理类
type AspectManager struct {
}

// Aspect 切面接口
type Aspect interface {
	Around(obj any, method *invoke.Method, args []reflect.Value)
}
