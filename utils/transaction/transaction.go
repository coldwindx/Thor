package transaction

import (
	"Thor/utils/invoke"
	"Thor/utils/proxy"
	"fmt"
	"reflect"
)

func init() {
	proxy.GlobalMethodProxyPublisher.AddSubscriber(&TransactionProxySubscriber{})
}

type TransactionProxySubscriber struct{}

func (t *TransactionProxySubscriber) Tag() string {
	return "transaction"
}

func (t *TransactionProxySubscriber) Validate(field reflect.StructField) (bool, string) {
	value, ok := field.Tag.Lookup(t.Tag())
	return ok, value
}

func (t *TransactionProxySubscriber) Proxy(proxy reflect.Value, field reflect.StructField, tag string) reflect.Value {
	proxy = reflect.MakeFunc(field.Type, func(args []reflect.Value) (results []reflect.Value) {
		invocation := &invoke.Method{Name: field.Name, Type: field.Type}
		fmt.Println("[transaction] >>>", field.Name, field.Type)
		results = invocation.Invoke(proxy.Interface(), args)
		return results
	})
	return proxy
}
