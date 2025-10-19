package invoke

import (
	"reflect"
)

// Method 表示一个方法的元数据
type Method struct {
	Name string            // 方法名
	Type reflect.Type      // 方法类型
	Tag  reflect.StructTag // 方法的tag
}

// Invoke 调用原始方法
// obj: 目标对象
// args: 方法参数
// 返回值: 方法返回值
func (im *Method) Invoke(obj any, args []reflect.Value) []reflect.Value {
	// 获取原始方法
	nativeMethod := reflect.ValueOf(obj).MethodByName(im.Name)
	// 检查方法是否存在
	if !nativeMethod.IsValid() {
		panic("method not found `" + im.Name + "` in `" + reflect.TypeOf(obj).String() + "`")
	}
	// 调用原始方法
	return nativeMethod.Call(args)
}

// InvocationMethod 表示一个方法调用的处理函数
type InvocationMethod func(obj any, method *Method, args []reflect.Value) []reflect.Value
