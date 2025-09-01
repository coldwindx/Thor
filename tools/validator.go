package tools

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateFuncStore 全局的验证方法集合
var ValidateFuncStore = make(map[string]validator.Func)

// 初始化
func init() {
	ValidateFuncStore["mobile"] = ValidateMobile
}

type Validator interface {
	GetValidateMessage() ValidatorMessages
}

type ValidatorMessages map[string]string

func GetErrorMsg(request interface{}, err error) string {
	var validationErrors validator.ValidationErrors
	isValidatorErrors := errors.As(err, &validationErrors)
	if !isValidatorErrors {
		return err.Error()
	}
	_, isValidator := request.(Validator)
	for _, v := range err.(validator.ValidationErrors) {
		if !isValidator {
			return v.Error()
		}
		message, exist := request.(Validator).GetValidateMessage()[v.Field()+"."+v.Tag()]
		if exist {
			return message
		}
	}
	return "Parameter error"
}

// 验证方法列表
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, mobile)
	return ok
}
