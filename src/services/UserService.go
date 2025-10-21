package services

import (
	"Thor/bootstrap"
	"Thor/bootstrap/inject"
	"Thor/src/models"
	"Thor/src/request"
	"Thor/utils"
	"errors"
	"gorm.io/gorm"
	"strconv"
)

type userService struct {
	Db *gorm.DB `inject:"MysqlConnection"`
}

var UserService = new(userService)

func init() {
	bootstrap.Beans.Provide(&inject.Object{Name: "UserService", Value: UserService})
}

func (t *userService) Register(params request.UserReq) (err error, user models.TUser) {
	var result = t.Db.Table("t_user").Where("mobile = ?", params.Mobile).Select("id").Find(&models.TUser{})
	if nil != result.Error {
		return errors.New(result.Error.Error()), user
	}
	if 0 != result.RowsAffected {
		return errors.New("手机号已存在"), user
	}
	user = models.TUser{
		Nickname: params.Name,
		Mobile:   params.Mobile,
		Email:    params.Email,
		Password: utils.BcryptMake([]byte(params.Password)),
	}
	err = t.Db.Table("t_user").Create(&user).Error
	return err, user
}

func (t *userService) Login(params request.LoginReq) (err error, user *models.TUser) {
	result := t.Db.Table("t_user").Where("mobile = ?", params.Mobile).First(&user)
	if nil != result.Error || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
		err = errors.New("用户名不存在或密码错误")
	}
	return err, user
}

func (t *userService) GetUserInfo(id string) (err error, user models.TUser) {
	iid, err := strconv.Atoi(id)
	result := t.Db.Table("t_user").First(&user, iid)
	return result.Error, user
}
