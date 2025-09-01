package request

import (
	"Thor/tools"
)

type UserReq struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// GetMessage 自定义错误信息
func (UserReq) GetValidateMessage() tools.ValidatorMessages {
	return tools.ValidatorMessages{
		"name.required":     "用户名称不能为空",
		"mobile.required":   "手机号码不能为空",
		"mobile.mobile":     "手机号码格式不正确",
		"email.required":    "邮箱不能为空",
		"password.required": "用户密码不能为空",
	}
}

type LoginReq struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (LoginReq) GetValidateMessage() tools.ValidatorMessages {
	return tools.ValidatorMessages{
		"mobile.required":   "手机号码不能为空",
		"mobile.mobile":     "手机号码格式不正确",
		"password.required": "用户密码不能为空",
	}
}
