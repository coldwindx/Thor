package bootstrap

import (
	"Thor/config"
	"Thor/utils/inject"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func init() {
	v := &RedisInitializer{name: "RedisInitializer", order: 3}
	Beans.Provide(&inject.Object{Name: v.GetName(), Value: v, Completed: true})
}

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
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Host + ":" + strconv.Itoa(config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	Beans.Provide(&inject.Object{Name: "RedisClient", Value: client, Completed: true})
}
func (ins *RedisInitializer) Close() {

}
