package bootstrap

import (
	"Thor/utils/inject"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"reflect"
	"strings"
)

func init() {
	v := &ValidatorInitializer{name: "ValidatorInitializer", order: 500}
	Beans.Provide(&inject.Object{Name: "ValidatorInitializer", Value: v, Completed: true})
}

type ValidatorInitializer struct {
	name       string
	order      int
	Validators map[string]Validator `inject:""`
}

func (v *ValidatorInitializer) GetName() string {
	return v.name
}

func (v *ValidatorInitializer) GetOrder() int {
	return v.order
}

func (v *ValidatorInitializer) Initialize() {
	// step1 启动验证器引擎
	engine := binding.Validator.Engine().(*validator.Validate)
	if nil == engine {
		panic(errors.New("validator_ engine not found"))
	}

	// step2 注册自定义验证器
	for _, validator_ := range v.Validators {
		if err := engine.RegisterValidation(validator_.GetName(), validator_.Validate); err != nil {
			panic(err)
		}
	}

	// step3 注册自定义json tag函数
	engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return lo.Ternary[string](name == "-", "", name)
	})
}

func (*ValidatorInitializer) Close() {

}

// Validator 自定义验证器接口
type Validator interface {
	GetName() string
	Validate(fl validator.FieldLevel) bool
}
