package logger

import (
	"fmt"
	"krathub/internal/conf"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	//配置zap日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	// 日志切割，采用 lumberjack 实现的
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "../../logs/" + c.Log.Filename, // 指定日志存储位置
		MaxSize:    10,                             // 日志的最大大小（M）
		MaxBackups: 5,                              // 日志的最大保存数量
		MaxAge:     30,                             // 日志文件存储最大天数
		Compress:   false,                          // 是否执行压缩
	})

	// 设置日志级别
	// TODO: 这里可以根据配置文件动态设置日志级别
	level := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	level.SetLevel(zapcore.Level(c.Log.Level))
	var core zapcore.Core

	// 根据配置文件的环境变量选择日志输出方式
	switch c.Env {
	case "dev":
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),                      // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台
			level, // 日志级别
		)
	default:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoder),                                      // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer), // 打印到控制台和文件
			level, // 日志级别
		)
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
