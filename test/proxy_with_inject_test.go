package test

import (
    "Thor/utils/inject"
    "Thor/utils/invoke"
    "fmt"
    "github.com/stretchr/testify/assert"
    "reflect"
    "testing"
)

var beans = &inject.Graph{}

func TestProxyBeforeInject(t *testing.T) {
    impl := &HelloWorld{Word: "Hello world"}

    proxy := invoke.NewMethodProxy(&Hello{}, func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
        return method.Invoke(impl, args)
    }).(*Hello)

    _ = beans.Provide(&inject.Object{Value: proxy})

    bean := beans.Query(reflect.TypeOf(&Hello{}))[0].Value.(*Hello)

    bean.SetWord("This is a proxy by HelloWorld")
    assert.Equal(t, "This is a proxy by HelloWorld", bean.SayHello())
    assert.Equal(t, "This is a proxy by HelloWorld", impl.SayHello())
}

// TestProxyAfterInject tests the proxy after the bean is injected.
func TestProxyAfterInject(t *testing.T) {
    impl := &HelloWorld{Word: "Hello world"}
    _ = beans.Provide(&inject.Object{Value: impl})

    proxy := invoke.NewMethodProxy(&Hello{}, func(obj any, method *invoke.Method, args []reflect.Value) []reflect.Value {
        fmt.Println("proxy method:", method.Name)
        return method.Invoke(impl, args)
    }).(*Hello)

    all := beans.Query(reflect.TypeOf(new(IHello)).Elem())
    if len(all) == 0 {
        t.Fatal("bean is nil")
    }
    bean := all[0].Value.(IHello)

    proxy.SetWord("This is a proxy by HelloWorld")

    t.Logf("start to test proxy after inject...")
    assert.Equal(t, "This is a proxy by HelloWorld", bean.SayHello())
    assert.Equal(t, "This is a proxy by HelloWorld", proxy.SayHello())
    assert.Equal(t, "This is a proxy by HelloWorld", impl.SayHello())
}
