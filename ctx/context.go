package ctx

import (
	"Thor/utils/inject"
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"github.com/zhuxiujia/GoMybatis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type MybatisMapperBind struct {
	XmlFile string `json:"xml_file"`
	Mapper  any    `json:"mapper"`
}

// 雪花算法ID
var Snowflake *snowflake.Node

var Statik http.FileSystem
var Logger *zap.Logger
var Db *gorm.DB

var MybatisEngine GoMybatis.GoMybatisEngine
var MybatisMapperBinds = make([]MybatisMapperBind, 0)
var DefaultSqlDB *sql.DB

// var MybatisMapperBinds = make([]func(), 0)
var Redis *redis.Client
var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

var Router *gin.Engine
var Routes = make([]func(*gin.Engine), 0)

// Bean 容器
var Beans = inject.NewGraph()
