package transaction

import (
	"Thor/utils/proxy"
	"context"
	"fmt"
	"reflect"
)

func init() {
	proxy.GlobalMethodProxyPublisher = append(proxy.GlobalMethodProxyPublisher, &TransactionProxySubscriber{})
}

type TransactionProxySubscriber struct{}

func (t *TransactionProxySubscriber) Before(tag string, obj any, args []reflect.Value) {
	if len(args) == 0 {
		return
	}
	// 从argos中获取第一个参数，并转换为context.Context类型
	ctx := args[0].Interface().(context.Context)
	// 给ctx添加一个transaction key
	ctx = context.WithValue(ctx, "transaction", tag)
	// 更新args[0]为新的context.Context对象
	args[0] = reflect.ValueOf(ctx)
	// 打印日志
	fmt.Println("[transaction] >>>", tag, obj, args)
}

func (t *TransactionProxySubscriber) After(tag string, obj any, args []reflect.Value, results []reflect.Value) {
	if len(args) == 0 || len(results) == 0 {
		return
	}
	// 从args中获取第一个参数，并转换为context.Context类型
	ctx := args[0].Interface().(context.Context)
	// 打印ctx中的transaction key
	fmt.Println("[transaction] ctx key >>>", ctx.Value("transaction"))
	// 从results中获取第一个参数，并转换为error类型
	err, ok := results[0].Interface().(error)
	if !ok {
		return
	}
	// 如果err不为空，说明事务执行过程中发生了错误
	if err != nil {
		// 打印日志
		fmt.Println("[transaction] <<<", tag, obj, args, results, err)
	} else {
		// 打印日志
		fmt.Println("[transaction] <<<", tag, obj, args, results)
	}
}

func (t *TransactionProxySubscriber) Tag() string {
	return "transaction"
}

func (t *TransactionProxySubscriber) Validate(field reflect.StructField) (bool, string) {
	value, ok := field.Tag.Lookup(t.Tag())
	return ok, value
}
