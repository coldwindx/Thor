package bootstrap

import (
	"Thor/config"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"strconv"
)

//func init() {
//	v := &RabbitMqInitializer{name: "rabbit", order: 100}
//	Initializers[v.name] = v
//}

type RabbitMqInitializer struct {
	name  string
	order int
}

func (ts *RabbitMqInitializer) GetName() string {
	return ts.name
}

func (ts *RabbitMqInitializer) GetOrder() int {
	return ts.order
}

func (ts *RabbitMqInitializer) Initialize() {
	log.Println("init rabbit...")
	var err error
	var c = config.Config.RabbitMq
	RabbitConn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/admin_vhost", c.User, c.Password, c.Addr, strconv.Itoa(c.Port)))
	if nil != err {
		Logger.Error("new mq conn err.", zap.Any("err", err))
		return
	}
	RabbitChannel, err = RabbitConn.Channel()
	if nil != err {
		Logger.Error("new mq conn err.", zap.Any("err", err))
	}
}

func (ts *RabbitMqInitializer) Close() {
	_ = RabbitChannel.Close()
	_ = RabbitConn.Close()
}
