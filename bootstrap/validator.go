package bootstrap

import (
	"Thor/ctx"
	"Thor/tools"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

func init() {
	v := &ValidatorInitializer{name: "validator", order: 500}
	Manager[v.name] = v
}

type ValidatorInitializer struct {
	name  string
	order int
}

func (t *ValidatorInitializer) GetName() string {
	return t.name
}

func (t *ValidatorInitializer) GetOrder() int {
	return t.order
}

func (*ValidatorInitializer) Initialize() {
	// step1 启动验证器引擎
	engine := binding.Validator.Engine()
	v, ok := engine.(*validator.Validate)
	if !ok {
		ctx.Logger.Error("load validator failed")
		return
	}

	// step2 注册自定义验证器
	for name, function := range tools.ValidateFuncStore {
		if err := v.RegisterValidation(name, function); nil != err {
			ctx.Logger.Error("load validate function fail, err:", zap.Any("err", err))
		}
	}
	// step3 注册自定义json tag函数
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if "-" == name {
			return ""
		}
		return name
	})
}

func (*ValidatorInitializer) Close() {

}
