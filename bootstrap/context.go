package bootstrap

import (
	"Thor/bootstrap/inject"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var Logger *zap.Logger

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

var Router *gin.Engine
var Routes = make([]func(*gin.Engine), 0)

// Bean 容器
var Beans = inject.NewGraph()
