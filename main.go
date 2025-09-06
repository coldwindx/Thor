//go:generate statik -src=./resources
//go:generate go fmt statik/statik.go

package main

import (
	"Thor/bootstrap"
	"Thor/ctx"
	_ "Thor/statik"
)

func main() {
	// step1 初始化配置
	bootstrap.Initialize()
	defer bootstrap.Close()
	ctx.Logger.Info("bootstrap init success!")
	// step 消息队列
	//go rabbitmq.Producer()
	//go rabbitmq.Consumer()
	// step 启动服务器
	bootstrap.Run()
}
