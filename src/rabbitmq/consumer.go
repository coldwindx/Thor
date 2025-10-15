package rabbitmq

import (
	"Thor/bootstrap"
	"fmt"
	"go.uber.org/zap"
)

func Consumer() {
	_ = bootstrap.RabbitChannel.Qos(1, 0, false)
	deliveries, err := bootstrap.RabbitChannel.Consume("default", "", false, false, false, false, nil)
	if nil != err {
		bootstrap.Logger.Error("从队列 default 获取数据失败", zap.Any("err", err))
	}

	for {
		select {
		case message := <-deliveries:
			body := string(message.Body)
			fmt.Println("消费数据：", body)
			//ctx.Logger.Info("消费数据：", zap.Any("body", body))
			_ = message.Ack(true)
		}
	}
}
