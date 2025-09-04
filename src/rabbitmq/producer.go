package rabbitmq

import (
	"Thor/ctx"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func Producer() {
	body := "hello"
	var msg = amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	err := ctx.RabbitChannel.Publish("default", "test", false, false, msg)
	if nil != err {
		ctx.Logger.Error("发送消息失败", zap.Any("err", err))
	}
}
