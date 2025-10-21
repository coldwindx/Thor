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

type JoinPoint struct {
	Obj    any             // 	目标对象
	Method *invoke.Method  //	目标方法
	Args   []reflect.Value //	目标方法的参数
}

// ProceedingJoinPoint 连接点结构体
type ProceedingJoinPoint struct {
	Index  int          //	责任链当前位置
	Points []*JoinPoint // 连接点列表
}

func (p *ProceedingJoinPoint) Append(point *JoinPoint) {
	p.Points = append(p.Points, point)
	p.Index = p.Index + 1
}

func (p *ProceedingJoinPoint) Proceed() []reflect.Value {
	// 调用代理方法
	point := p.Points[p.Index-1]
	p.Index = p.Index - 1
	// 调用下一个连接点
	result := point.Method.Invoke(point.Obj, point.Args)
	return result
}

func (p *ProceedingJoinPoint) GetObject() any {
	return p.Points[0].Obj
}

func (p *ProceedingJoinPoint) GetMethod() *invoke.Method {
	return p.Points[0].Method
}

func (p *ProceedingJoinPoint) GetArgs() []reflect.Value {
	return p.Points[0].Args
}

// AspectManager 统一的切面管理类
type AspectManager struct {
	Aspects []Aspect `inject:""`
}

// Aspect 切面接口
type Aspect interface {
	// Pointcut 切点
	Pointcut(method *invoke.Method) bool
	// Around 环绕通知
	Around(jcp *ProceedingJoinPoint) []reflect.Value
}
