package aop

import (
	"Thor/bootstrap"
	"Thor/utils/inject"
	"Thor/utils/invoke"
	"reflect"
)

func init() {
	bootstrap.Beans.Provide(&inject.Object{Name: "AspectManager", Value: &AspectManager{}})
}

// ProceedingJoinPoint 连接点结构体
type ProceedingJoinPoint struct {
	Obj    any             // 目标对象
	Method *invoke.Method  //	目标方法
	Args   []reflect.Value //	目标方法的参数
}

func (p *ProceedingJoinPoint) Proceed() []reflect.Value {
	// 调用目标方法
	return p.Method.Invoke(p.Obj, p.Args)
}

// AspectManager 统一的切面管理类
type AspectManager struct {
	Aspects []Aspect `inject:""`
}

// Aspect 切面接口
type Aspect interface {
	// Around 环绕通知
	Around(jcp *ProceedingJoinPoint) []reflect.Value
}
