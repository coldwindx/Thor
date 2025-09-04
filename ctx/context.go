package ctx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Logger *zap.Logger
var Db *gorm.DB
var Redis *redis.Client
var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel
var Router *gin.Engine

var Routes = make([]func(*gin.Engine), 0)
