package main

import (
	"Thor/bootstrap"
	"Thor/ctx"
)

func main() {
	// step1 初始化配置
	bootstrap.Initialize()
	defer bootstrap.Close()
	ctx.Logger.Info("bootstrap init success!")
	// step 启动服务器
	bootstrap.Run()
}
