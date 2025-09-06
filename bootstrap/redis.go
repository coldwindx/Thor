package bootstrap

import (
	"Thor/config"
	"Thor/ctx"
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
)

//func init() {
//	v := &RedisInitializer{name: "redis", order: 3}
//	Manager[v.name] = v
//}

type RedisInitializer struct {
	name  string
	order int
}

func (ins *RedisInitializer) GetName() string {
	return ins.name
}
func (ins *RedisInitializer) GetOrder() int {
	return ins.order
}
func (ins *RedisInitializer) Initialize() {
	ctx.Redis = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Host + ":" + strconv.Itoa(config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})

	_, err := ctx.Redis.Ping(context.Background()).Result()
	if nil != err {
		ctx.Logger.Error("redis connect ping failed.", zap.Any("err", err))
	}
}
func (ins *RedisInitializer) Close() {

}
