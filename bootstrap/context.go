package bootstrap

import (
	"Thor/utils/inject"
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/zhuxiujia/GoMybatis"
	"go.uber.org/zap"
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

var MybatisEngine GoMybatis.GoMybatisEngine
var MybatisMapperBinds = make([]MybatisMapperBind, 0)
var DefaultSqlDB *sql.DB

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

var Router *gin.Engine
var Routes = make([]func(*gin.Engine), 0)

// Bean 容器
var Beans = inject.NewGraph()
