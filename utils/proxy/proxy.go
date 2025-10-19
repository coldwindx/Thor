package proxy

import (
	"Thor/utils/invoke"
	"reflect"
)

var GlobalMethodProxyPublisher = make([]MethodProxyWrapper, 0)

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
	// DEBUG 打印itf的类型
	//fmt.Println("[itf] >>>", itfType)
	// 遍历指针对象的所有方法
	for i := 0; i < reflect.ValueOf(itf).Elem().NumField(); i++ {
		field := itfType.Elem().Field(i)
		value := itfValue.Elem().Field(i)
		//fmt.Println("[field] >>>", field.Name, field.Type)
		// 检查方法是否有效
		if field.Type.Kind() != reflect.Func || !value.CanSet() {
			continue
		}

		//fmt.Println("[proxy] >>>", field.Name, field.Type)
		// 创建代理方法
		proxy := reflect.MakeFunc(field.Type, func(args []reflect.Value) []reflect.Value {
			// 获取可行代理
			wrappers := make([]MethodProxyPackage, 0)
			for _, wrapper := range GlobalMethodProxyPublisher {
				if valid, tag := wrapper.Validate(field); valid {
					w := MethodProxyPackage{TagName: wrapper.Tag(), TagValue: tag, Wrapper: wrapper}
					wrappers = append(wrappers, w)
				}
			}
			// 前置代理
			for _, w := range wrappers {
				w.Wrapper.Before(w.TagValue, itf, args)
			}
			// 调用原始方法
			invocation := &invoke.Method{Name: field.Name, Type: field.Type}
			results := handler(itf, invocation, args)
			// 后置代理
			for _, w := range wrappers {
				w.Wrapper.After(w.TagValue, itf, args, results)
			}
			return results
		})
		value.Set(proxy)
	}
	return itf
}

// MethodProxyWrapper 具体的方法代理订阅者
type MethodProxyWrapper interface {
	Tag() string
	Validate(field reflect.StructField) (bool, string)
	Before(string, any, []reflect.Value)
	After(string, any, []reflect.Value, []reflect.Value)
}

type MethodProxyPackage struct {
	TagName  string
	TagValue string
	Wrapper  MethodProxyWrapper
}
