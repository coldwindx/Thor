package main

import (
	"Thor/bootstrap"
	"Thor/ctx"
	"Thor/src/rabbitmq"
)

func main() {
	// step1 初始化配置
	bootstrap.Initialize()
	defer bootstrap.Close()
	ctx.Logger.Info("bootstrap init success!")
	// step 消息队列
	go rabbitmq.Producer()
	go rabbitmq.Consumer()
	// step 启动服务器
	bootstrap.Run()
}
