package initializer

import (
	"Thor/bootstrap"
	"Thor/bootstrap/inject"
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func init() {
	v := &RedisInitializer{name: "RedisInitializer", order: 3}
	bootstrap.Beans.Provide(&inject.Object{Name: v.GetName(), Value: v, Completed: true})
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
		Addr:     bootstrap.Config.Redis.Host + ":" + strconv.Itoa(bootstrap.Config.Redis.Port),
		Password: bootstrap.Config.Redis.Password,
		DB:       bootstrap.Config.Redis.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	bootstrap.Beans.Provide(&inject.Object{Name: "RedisClient", Value: client, Completed: true})
}
func (ins *RedisInitializer) Close() {

}
