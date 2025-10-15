package bootstrap

import (
	"Thor/config"
	"Thor/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	level   zapcore.Level
	options []zap.Option
)

func init() {
	ins := &LogInitializer{name: "log", order: 2}
	Manager[ins.name] = ins
}

type LogInitializer struct {
	name  string
	order int
}

func (t *LogInitializer) GetName() string {
	return t.name
}
func (t *LogInitializer) GetOrder() int {
	return t.order
}
func (t *LogInitializer) Initialize() {
	// step1 创建根目录
	if ok, _ := utils.PathExists(config.Config.Log.Dir); !ok {
		_ = os.Mkdir(config.Config.Log.Dir, os.ModePerm)
	}
	// step2 设置日志等级
	switch config.Config.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if config.Config.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}
	// step3 初始化zap
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(l.String())
	}

	if "json" == config.Config.Log.Format {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder, getLogWriter(), level)
	Logger = zap.New(core, options...)
}

func (t *LogInitializer) Close() {

}
func getLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   config.Config.Log.Dir + "/" + config.Config.Log.Filename,
		MaxSize:    config.Config.Log.MaxSize,
		MaxBackups: config.Config.Log.MaxBackups,
		MaxAge:     config.Config.Log.MaxAge,
		Compress:   config.Config.Log.Compress,
	}
	return zapcore.AddSync(file)
}
