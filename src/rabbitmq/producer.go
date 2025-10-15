package rabbitmq

import (
	"Thor/bootstrap"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func Producer() {
	body := "hello"
	var msg = amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	err := bootstrap.RabbitChannel.Publish("default", "test", false, false, msg)
	if nil != err {
		bootstrap.Logger.Error("发送消息失败", zap.Any("err", err))
	}
}
