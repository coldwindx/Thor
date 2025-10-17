package transaction

import (
    "Thor/utils/proxy"
    "fmt"
    "reflect"
)

func init() {
    proxy.GlobalMethodProxyPublisher = append(proxy.GlobalMethodProxyPublisher, &TransactionProxySubscriber{})
}

type TransactionProxySubscriber struct{}

func (t *TransactionProxySubscriber) Before(tag string, obj any, args []reflect.Value) {
    fmt.Println("[transaction] >>>", tag, obj, args)
}

func (t *TransactionProxySubscriber) After(tag string, obj any, args []reflect.Value, results []reflect.Value) {
    fmt.Println("[transaction] <<<", tag, obj, args, results)
}

func (t *TransactionProxySubscriber) Tag() string {
    return "transaction"
}

func (t *TransactionProxySubscriber) Validate(field reflect.StructField) (bool, string) {
    value, ok := field.Tag.Lookup(t.Tag())
    return ok, value
}
