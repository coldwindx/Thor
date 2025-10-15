package inject

import (
	"Thor/utils/invoke"
	"github.com/samber/lo"
	"reflect"
	"strings"
)

// CycleProvide 循环提供Bean实例，需要制定bean tag标签
func (g *Graph) CycleProvide(objs ...*Object) {
	// 遍历objs数组
	for _, obj := range objs {
		obj.RfType = reflect.TypeOf(obj.Value)
		obj.RfValue = reflect.ValueOf(obj.Value)
		// 遍历obj的所有属性值
		for i := 0; i < obj.RfType.Elem().NumField(); i++ {
			// 获取属性的类型
			field := obj.RfType.Elem().Field(i)
			fieldValue := obj.RfValue.Elem().Field(i)
			// 跳过不能设置的属性，或者属性不是对象指针类型
			if !fieldValue.CanSet() || field.Type.Kind() != reflect.Ptr || field.Type.Elem().Kind() != reflect.Struct {
				continue
			}
			// 获取这个属性的tag标签
			tag, ok := field.Tag.Lookup("bean")
			if !ok {
				continue
			}
			// 根据;分割字符串，获取 bean的名称 和 bean的代理模式
			tags := strings.Split(tag, ";")
			// 如果没有指定bean tag的名称，则默认使用属性的类型名称
			name := lo.Ternary(len(tags[0]) == 0, field.Type.Elem().Name(), tags[0])
			// 根据tag标签，从容器中获取对应的bean对象，如果存在，说明不需要自动创建
			if _, ok = g.objects[name]; ok {
				continue
			}
			// 需要自动创建bean对象
			bean := reflect.New(field.Type.Elem()).Interface()
			g.Provide(&Object{Name: name, Value: bean})
			// 自动代理
			if len(tags) <= 1 || tags[1] != "proxy" {
				continue
			}
			invoke.NewMethodProxy(bean, func(_ any, method *invoke.Method, args []reflect.Value) []reflect.Value {
				return method.Invoke(obj.Value, args)
			})
		}

		// 检查Bean实例是否是指针类型
		if obj.RfType.Kind() != reflect.Ptr {
			panic("Bean实例必须是指针类型")
		}
		// 检查Bean的名称是否为空，如果为空，则根据bean的类型构造名称
		obj.Name = lo.Ternary(len(obj.Name) == 0, obj.RfType.Elem().Name(), obj.Name)
		// 检查是否有相同名称的bean对象
		if _, ok := g.objects[obj.Name]; ok {
			panic("provided two instances named `" + obj.Name + "`")
		}
		// 根据Bean的名称，加入容器管理
		g.objects[obj.Name] = obj
	}
}
