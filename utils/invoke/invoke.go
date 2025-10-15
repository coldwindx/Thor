package invoke

import (
    "reflect"
)

// Method 表示一个方法的元数据
type Method struct {
    Name string       // 方法名
    Type reflect.Type // 方法类型
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

// NewMethodProxy 创建一个代理对象
// itf: 接口指针
// handler: 方法调用处理函数
// 返回值: 代理对象
func NewMethodProxy(itf any, handler InvocationMethod) any {
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
        // 代理方法
        //fmt.Println("[proxy] >>>", field.Name, field.Type)
        proxy := reflect.MakeFunc(field.Type, func(args []reflect.Value) []reflect.Value {
            invocation := &Method{Name: field.Name, Type: field.Type}
            return handler(itf, invocation, args)
        })
        value.Set(proxy)
    }
    return itf
}
