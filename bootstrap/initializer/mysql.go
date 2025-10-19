package initializer

import (
	"Thor/bootstrap"
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
	bootstrap.Beans.Provide(&inject.Object{Name: initializer.GetName(), Value: initializer, Completed: true})
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
	dbConfig := bootstrap.Config.Database
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
		bootstrap.Logger.Error("mysql connect failed, err:", zap.Any("err", err))
		return
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// 管理数据库连接
	bootstrap.Beans.Provide(&inject.Object{Name: "MysqlConnection", Value: db, Completed: true})
	bootstrap.Beans.Provide(&inject.Object{Name: "DBClient", Value: &DBClient{db: db}, Completed: true})
}

func (t *DatabaseInitializer) Close() {
	db := bootstrap.Beans.GetByName("MysqlConnection").(*gorm.DB)
	if db != nil {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}
}

// 自定义gorm writer
func getGormLogWriter() logger.Writer {
	var writer io.Writer = os.Stdout
	if bootstrap.Config.Database.EnableFileLogWriter {
		writer = &lumberjack.Logger{
			Filename:   bootstrap.Config.Log.Dir + "/" + bootstrap.Config.Database.LogFilename,
			MaxSize:    bootstrap.Config.Log.MaxSize,
			MaxBackups: bootstrap.Config.Log.MaxBackups,
			MaxAge:     bootstrap.Config.Log.MaxAge,
			Compress:   bootstrap.Config.Log.Compress,
		}
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel
	switch bootstrap.Config.Database.LogMode {
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
		Colorful:                  !bootstrap.Config.Database.EnableFileLogWriter,
	})
}
