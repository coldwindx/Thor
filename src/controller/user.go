package controller

import (
	"Thor/common"
	"Thor/ctx"
	"Thor/src/middleware"
	"Thor/src/models"
	"Thor/src/request"
	"Thor/src/services"
	"Thor/tools"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("insert route")
	ctx.Routes = append(ctx.Routes, func(r *gin.Engine) {
		group1 := r.Group("/user")
		group1.POST("/register", Register)

		group2 := r.Group("/auth")
		group2.POST("/login", Login)

		group2WithAuth := group2.Use(middleware.JWTAuth(services.AppGuardName))
		group2WithAuth.POST("/info", Info)
		group2WithAuth.POST("/logout", Logout)
	})
}

func Register(c *gin.Context) {
	var form request.UserReq
	if err := c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}
	err, user := services.UserService.Register(form)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
	} else {
		common.Success(c, user)
	}
}

func Login(c *gin.Context) {
	var err error
	var form request.LoginReq
	if err = c.ShouldBindJSON(&form); nil != err {
		common.Fail(c, common.Errors.ValidateError.Code, tools.GetErrorMsg(form, err))
		return
	}

	var user *models.TUser = nil
	if err, user = services.UserService.Login(form); nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	token, err, _ := services.JwtService.CreateToken(services.AppGuardName, user)
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, token)
}

func Logout(c *gin.Context) {
	err := services.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, "登出失败："+err.Error())
		return
	}
	common.Success(c, nil)
}

func Info(c *gin.Context) {
	err, user := services.UserService.GetUserInfo(c.Keys["id"].(string))
	if nil != err {
		common.Fail(c, common.Errors.BusinessError.Code, err.Error())
		return
	}
	common.Success(c, user)
}
