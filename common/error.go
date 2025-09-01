package common

type Error struct {
	Code    int
	Message string
}

type CustomErrors struct {
	BusinessError Error
	TokenError    Error
	ValidateError Error
}

var Errors = CustomErrors{
	BusinessError: Error{40000, "业务错误"},
	TokenError:    Error{40100, "登录鉴权失效"},
	ValidateError: Error{42200, "请求参数错误"},
}
