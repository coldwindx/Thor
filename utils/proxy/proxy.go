package proxy

import (
	"Thor/utils/invoke"
	"reflect"
)

var GlobalMethodProxyPublisher = &MethodProxyPublisher{
	subscribers: make(map[string]MethodProxySubscriber),
}

// NewMethodProxy 创建一个代理对象
// itf: 接口指针
// handler: 方法调用处理函数
// 返回值: 代理对象
func NewMethodProxy(itf any, handler invoke.InvocationMethod) any {
	// 获取接口类型
	itfType := reflect.TypeOf(itf)
	itfValue := reflect.ValueOf(itf)
	// 类型校验
	if itfType.Kind() != reflect.Ptr || itfType.Elem().Kind() != reflect.Struct {
		panic("Need a pointer of struct")
	}

	// 遍历指针对象的所有方法
	for i := 0; i < reflect.ValueOf(itf).Elem().NumField(); i++ {
		field := itfType.Elem().Field(i)
		value := itfValue.Elem().Field(i)
		//fmt.Println("[field] >>>", field.Name, field.Type)
		// 检查方法是否有效
		if field.Type.Kind() != reflect.Func || !value.CanSet() {
			continue
		}
		// 先代理底层方法
		//fmt.Println("[proxy] >>>", field.Name, field.Type)
		proxy := reflect.MakeFunc(field.Type, func(args []reflect.Value) []reflect.Value {
			invocation := &invoke.Method{Name: field.Name, Type: field.Type}
			return handler(itf, invocation, args)
		})
		// 再根据tag判断是否需要额外代理
		proxy = GlobalMethodProxyPublisher.Publish(proxy, field)
		value.Set(proxy)
	}
	return itf
}

// MethodProxySubscriber 具体的方法代理订阅者
type MethodProxySubscriber interface {
	Tag() string
	Validate(field reflect.StructField) (bool, string)
	Proxy(proxy reflect.Value, field reflect.StructField, tag string) reflect.Value
}
type MethodProxyPublisher struct {
	subscribers map[string]MethodProxySubscriber
}

func (p *MethodProxyPublisher) AddSubscriber(subscriber MethodProxySubscriber) {
	p.subscribers[subscriber.Tag()] = subscriber
}

func (p *MethodProxyPublisher) Publish(proxy reflect.Value, field reflect.StructField) reflect.Value {
	for _, subscriber := range p.subscribers {
		valid, tag := subscriber.Validate(field)
		if !valid {
			continue
		}
		// 附加代理
		proxy = subscriber.Proxy(proxy, field, tag)
	}
	return proxy
}
