package services

import (
	"Thor/bootstrap"
	"Thor/bootstrap/inject"
	"Thor/utils"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type jwtService struct {
	Redis *redis.Client `inject:"RedisClient"`
}

var JwtService = new(jwtService)

func init() {
	bootstrap.Beans.Provide(&inject.Object{Name: "JwtService", Value: JwtService})
}

// 所有需要token的用户模型接口
type JwtUser interface {
	GetUid() string
}

type CustomClaims struct {
	jwt.StandardClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (service *jwtService) CreateToken(GuardName string, user JwtUser) (TokenOutPut, error, *jwt.Token) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + bootstrap.Config.Jwt.JwtTtl,
		Id:        user.GetUid(),
		Issuer:    GuardName, // 用于在中间件中区分不同客户端颁发的token
		NotBefore: time.Now().Unix() - 1000,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{StandardClaims: claims})
	tokenStr, err := token.SignedString([]byte(bootstrap.Config.Jwt.Secret))
	return TokenOutPut{tokenStr, int(bootstrap.Config.Jwt.JwtTtl), TokenType}, err, token
}

func (service *jwtService) getBlackListKey(token string) string {
	return "jwt_black_list:" + utils.MD5([]byte(token))
}

func (service *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	now := time.Now().Unix()
	timer := time.Duration(token.Claims.(*CustomClaims).ExpiresAt-now) * time.Second
	result := service.Redis.SetNX(context.Background(), service.getBlackListKey(token.Raw), now, timer)
	return result.Err()
}

func (service *jwtService) IsInBlackList(token string) bool {
	key := service.getBlackListKey(token)
	joinUnixStr, err := service.Redis.Get(context.Background(), key).Result()
	if nil != err {
		bootstrap.Logger.Error("get key from redis fail.", zap.Any("key", key), zap.Any("err", err))
		return false
	}
	if "" == joinUnixStr {
		bootstrap.Logger.Error("value got from redis is empty.", zap.Any("key", key))
		return false
	}

	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if nil != err {
		bootstrap.Logger.Error("strconv value fail.", zap.Any("value", joinUnixStr), zap.Any("err", err))
		return false
	}
	return bootstrap.Config.Jwt.JwtBlacklistGracePeriod <= time.Now().Unix()-joinUnix
}

func (service *jwtService) GetUserInfo(GuardName string, id string) (err error, user JwtUser) {
	switch GuardName {
	case AppGuardName:
		return UserService.GetUserInfo(id)
	default:
		err = errors.New("guard " + GuardName + " does not exist")
	}
	return err, nil
}
