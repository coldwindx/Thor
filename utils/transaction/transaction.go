package transaction

import (
	"Thor/bootstrap"
	"Thor/utils/aop"
	"Thor/utils/inject"
	"context"
	"errors"
	"github.com/samber/lo"
	"reflect"
)

func init() {
	bootstrap.Beans.Provide(&inject.Object{Name: "aop.Transaction", Value: &TransactionAspect{}})
}

// TransactionAspect 事务切面
type TransactionAspect struct {
}

func (t TransactionAspect) Around(jcp *aop.ProceedingJoinPoint) []reflect.Value {
	client := bootstrap.Beans.GetByName("DBClient").(*bootstrap.DBClient)
	// 从切面中获取方法请求参数
	args := jcp.Args
	if 0 == len(args) {
		panic("first argument must be context.Context when use transaction aspect")
	}
	ctx, ok := jcp.Args[0].Interface().(context.Context)
	if !ok {
		panic("first argument must be context.Context when use transaction aspect")
	}
	// 开启事务
	var res []reflect.Value
	if err := client.Transaction(ctx, func(ctx context.Context) error {
		// 执行目标方法
		res = jcp.Proceed()
		// 检查目标方法是否返回错误
		if 0 == len(res) {
			return errors.New("must return at least one error when use transaction aspect")
		}
		// 检查目标方法是否返回错误
		errVal, ok := res[len(res)-1].Interface().(error)
		return lo.Ternary(ok, errVal, errors.New("must return error at last when use transaction aspect"))
	}); err != nil {
		panic(err)
	}
	// 返回目标方法的执行结果
	return res
}
