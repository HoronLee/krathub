package logger

import (
	"fmt"
	"krathub/internal/conf"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_ log.Logger = (*ZapLogger)(nil)
)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// NewLogger 配置zap日志,将zap日志库引入
func NewLogger(c *conf.App) log.Logger {
	// 设置日志级别
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if c.Log != nil {
		level.SetLevel(zapcore.Level(c.Log.Level))
	}

	// lumberjack 日志切割
	var lumberjackLogger *lumberjack.Logger
	if c.Log != nil {
		maxSize := 10
		if c.Log.MaxSize != 0 {
			maxSize = int(c.Log.MaxSize)
		}
		maxBackups := 5
		if c.Log.MaxBackups != 0 {
			maxBackups = int(c.Log.MaxBackups)
		}
		maxAge := 30
		if c.Log.MaxAge != 0 {
			maxAge = int(c.Log.MaxAge)
		}
		lumberjackLogger = &lumberjack.Logger{
			Filename:   "../../logs/" + c.Log.Filename, // 日志文件路径
			MaxSize:    maxSize,                        // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: maxBackups,                     // 日志文件最多保存多少个备份
			MaxAge:     maxAge,                         // 文件最多保存多少天
			Compress:   c.Log.Compress,                 // 是否压缩
		}
	}

	// 根据不同环境设置不同的日志输出
	var core zapcore.Core
	switch c.Env {
	case "dev":
		// dev模式，终端彩色输出，不输出到文件
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 显式设置彩色日志级别
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	case "prod":
		// prod模式，终端非json非彩色输出，文件json非彩色输出
		// 可以采用Unix timeStamp或ISO8601时间格式
		prodEncoderConfig := zap.NewProductionEncoderConfig()
		prodEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(prodEncoderConfig)
		jsonEncoder := zapcore.NewJSONEncoder(prodEncoderConfig)
		if lumberjackLogger == nil {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
				zapcore.NewCore(jsonEncoder, zapcore.AddSync(lumberjackLogger), level),
			)
		}
	case "test":
		// test模式，不输出日志
		core = zapcore.NewNopCore()
	default:
		// 默认情况，使用prod模式
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		if lumberjackLogger == nil {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
				zapcore.NewCore(jsonEncoder, zapcore.AddSync(lumberjackLogger), level),
			)
		}
	}

	opts := []zap.Option{
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	}

	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// Log 实现log接口
func (l *ZapLogger) Log(level log.Level, keyvals ...any) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	case log.LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}

func (l *ZapLogger) GetGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     l.log,
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值，单位为千分之一秒
	}
}
