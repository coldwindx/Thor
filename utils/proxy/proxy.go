package proxy

import (
    "Thor/bootstrap"
    "Thor/utils/aop"
    "Thor/utils/inject"
    "Thor/utils/invoke"
    "github.com/samber/lo"
    "reflect"
    "strings"
)

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
            // ******************************* 切面责任链 ************************************** //
            // 获取切面管理器
            aspectManager := bootstrap.Beans.GetByName("AspectManager").(*aop.AspectManager)
            method := &invoke.Method{Name: field.Name, Type: field.Type}
            chain := &aop.ProceedingJoinPoint{Obj: itf, Method: method, Args: args}
            // 构建责任链
            aopMethod := &invoke.Method{Name: "Around", Type: reflect.TypeOf(aop.Aspect.Around)}
            for _, aspect := range aspectManager.Aspects {
                // 构建责任链
                args = []reflect.Value{reflect.ValueOf(chain)}
                chain = &aop.ProceedingJoinPoint{Obj: aspect, Method: aopMethod, Args: args}
            }
            // ******************************* 切面责任链 ************************************** //
            // 调用责任链
            return chain.Proceed()
            //invocation := &invoke.Method{Name: field.Name, Type: field.Type}
            //results := handler(itf, invocation, args)
            //return results
        })
        value.Set(proxy)
    }
    return itf
}

// CycleProvide 循环提供Bean实例，需要制定bean tag标签
func CycleProvide(g *inject.Graph, objs ...*inject.Object) {
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
            if ok = g.ExistByName(name); ok {
                continue
            }
            // 需要自动创建bean对象
            bean := reflect.New(field.Type.Elem()).Interface()
            g.Provide(&inject.Object{Name: name, Value: bean})
            // 自动代理
            if len(tags) <= 1 || tags[1] != "proxy" {
                continue
            }
            NewMethodProxy(bean, func(_ any, method *invoke.Method, args []reflect.Value) []reflect.Value {
                return method.Invoke(obj.Value, args)
            })
        }
    }
    // 注入原始bean
    g.Provide(objs...)
}
