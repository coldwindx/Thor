package test

import (
    "Thor/utils/invoke"
    "fmt"
    "reflect"
    "testing"
)

func TestNewInstance(t *testing.T) {
    proxy := invoke.NewMethodProxy(&Hello{}, func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
        return []reflect.Value{reflect.ValueOf("This is a proxy function")}
    }).(*Hello)
    result := proxy.SayHello()
    fmt.Println(result)
}

func TestAOPProxyWithClass(t *testing.T) {
    impl := &HelloWorld{Word: "Hello world"}

    proxy := invoke.NewMethodProxy(&Hello{}, func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
        return method.Invoke(impl, args)
    }).(*Hello)

    proxy.SetWord("This is a proxy by HelloWorld")
    fmt.Println(proxy.SayHello())
    fmt.Println(impl.SayHello())
}

func TestAOPProxyWithSameStruct(t *testing.T) {
    impl := &HelloWorld{Word: "Hello world"}

    proxy := invoke.NewMethodProxy(&HelloWorld{}, func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
        return method.Invoke(impl, args)
    }).(*HelloWorld)

    proxy.SetWord("This is a proxy by HelloWorld")
    fmt.Println(proxy.SayHello())
    fmt.Println(impl.SayHello())
}

// 不支持接口代理
//func TestAOPProxyWithInterface(t *testing.T) {
//    impl := &HelloWorld{Word: "Hello world"}
//
//    proxy := invoke.NewMethodProxy(new(IHello), func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
//        return method.Invoke(impl, args)
//    }).(IHello)
//
//    proxy.SetWord("This is a proxy by HelloWorld")
//    fmt.Println(proxy.SayHello())
//    fmt.Println(impl.SayHello())
//}
