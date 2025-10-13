package test

import (
	"Thor/utils/invoke"
	"fmt"
	"reflect"
)

func init() {
	fmt.Println("Testing proxy...")
	testProxyInImpl()
}

type Hello interface {
	SayHello() string
	SetWord(word string)
}
type HelloWorld struct {
	Word string
}

func (h *HelloWorld) SayHello() string {
	return h.Word
}

func (h *HelloWorld) SetWord(word string) {
	h.Word = word
}

type SimpleBeanFactory struct {
}

func (sbf SimpleBeanFactory) NewInstance(itf any) any {
	proxy := invoke.NewMethodProxy(itf, func(obj any, method invoke.Method, args []reflect.Value) []reflect.Value {
		fmt.Println("invoke method", method)
		result := method.Invoke(obj, args)
		fmt.Println("invoke result", result)
		return result
	})

	return proxy
}

func testProxyInImpl() {
	myService := &HelloWorld{Word: "hello"}
	// 创建动态代理对象
	sbf := SimpleBeanFactory{}
	proxy := sbf.NewInstance(myService)
	// 调用代理对象的方法
	fmt.Println(proxy.(*HelloWorld).SayHello())
}
