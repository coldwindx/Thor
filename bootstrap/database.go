package bootstrap

import (
	"Thor/config"
	"Thor/utils/inject"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	initializer := &DatabaseInitializer{name: "DatabaseInitializer", order: 100}
	Beans.Provide(&inject.Object{Name: initializer.GetName(), Value: initializer, Completed: true})
}

type DatabaseInitializer struct {
	name  string
	order int
}

func (t *DatabaseInitializer) GetName() string {
	return t.name
}
func (t *DatabaseInitializer) GetOrder() int {
	return t.order
}
func (t *DatabaseInitializer) Initialize() {
	dbConfig := config.Config.Database
	if "" == dbConfig.Database {
		return
	}
	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameColumn:   true,
		DontSupportRenameIndex:    true,
		SkipInitializeWithVersion: false,
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   getGormLogger(),
	})
	if nil != err {
		Logger.Error("mysql connect failed, err:", zap.Any("err", err))
		return
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// 管理数据库连接
	Beans.Provide(&inject.Object{Name: "MysqlConnection", Value: db, Completed: true})
}

func (t *DatabaseInitializer) Close() {
	db := Beans.GetByName("MysqlConnection").(*gorm.DB)
	if db != nil {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}
}

// 自定义gorm writer
func getGormLogWriter() logger.Writer {
	var writer io.Writer = os.Stdout
	if config.Config.Database.EnableFileLogWriter {
		writer = &lumberjack.Logger{
			Filename:   config.Config.Log.Dir + "/" + config.Config.Database.LogFilename,
			MaxSize:    config.Config.Log.MaxSize,
			MaxBackups: config.Config.Log.MaxBackups,
			MaxAge:     config.Config.Log.MaxAge,
			Compress:   config.Config.Log.Compress,
		}
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel
	switch config.Config.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}
	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Microsecond,
		LogLevel:                  logMode,
		IgnoreRecordNotFoundError: true,
		Colorful:                  !config.Config.Database.EnableFileLogWriter,
	})
}
