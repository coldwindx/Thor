package bootstrap

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

func init() {
	v := &SnowflakeInitializer{name: "snowflake", order: 10}
	Manager[v.name] = v
}

type SnowflakeInitializer struct {
	name  string
	order int
}

func (it *SnowflakeInitializer) GetName() string {
	return it.name
}

func (it *SnowflakeInitializer) GetOrder() int {
	return it.order
}

func (it *SnowflakeInitializer) Initialize() {
	var err error

	parse, err := time.Parse("2006-01-02", "2025-09-07")
	if err != nil {
		panic("雪花算法ID构造失败，初始化时间错误." + err.Error())
	}
	snowflake.Epoch = parse.UnixNano() / 1e6

	Snowflake, err = snowflake.NewNode(1)
	if err != nil {
		panic("雪花算法ID构造失败." + err.Error())
	}
}

func (it *SnowflakeInitializer) Close() {

}
