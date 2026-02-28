package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ log.Logger = (*ZapLogger)(nil)

type Logger = log.Logger
type Helper = log.Helper

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

type Config struct {
	Env        string
	Level      int32
	Filename   string
	MaxSize    int32
	MaxBackups int32
	MaxAge     int32
	Compress   bool
}

// NewLogger 配置zap日志,将zap日志库引入
func NewLogger(c *Config) log.Logger {
	if c != nil && c.Filename == "" {
		c.Filename = "./logs/app.log"
	}

	// 设置日志级别，只支持 Kratos 定义的 5 个等级
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if c != nil {
		// 映射 Kratos 日志等级到 Zap 等级，并对无效值进行默认处理
		switch c.Level {
		case 0: // debug
			level.SetLevel(zapcore.DebugLevel)
		case 1: // info
			level.SetLevel(zapcore.InfoLevel)
		case 2: // warn
			level.SetLevel(zapcore.WarnLevel)
		case 3: // error
			level.SetLevel(zapcore.ErrorLevel)
		case 4: // fatal
			level.SetLevel(zapcore.FatalLevel)
		default:
			// 无效等级默认使用 info
			level.SetLevel(zapcore.InfoLevel)
		}
	}

	// lumberjack 日志切割
	var lumberjackLogger *lumberjack.Logger
	if c != nil {
		maxSize := 10
		if c.MaxSize != 0 {
			maxSize = int(c.MaxSize)
		}
		maxBackups := 5
		if c.MaxBackups != 0 {
			maxBackups = int(c.MaxBackups)
		}
		maxAge := 30
		if c.MaxAge != 0 {
			maxAge = int(c.MaxAge)
		}
		if dir := filepath.Dir(c.Filename); dir != "." && dir != "/" {
			_ = os.MkdirAll(dir, 0755)
		}
		lumberjackLogger = &lumberjack.Logger{
			Filename:   c.Filename, // 日志文件路径
			MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: maxBackups, // 日志文件最多保存多少个备份
			MaxAge:     maxAge,     // 文件最多保存多少天
			Compress:   c.Compress, // 是否压缩
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

func NewHelper(logger Logger, opts ...Option) *Helper {
	if len(opts) > 0 {
		logger = With(logger, opts...)
	}
	return log.NewHelper(logger)
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

// GetGormLogger 获取Gorm日志适配器
func (l *ZapLogger) GetGormLogger(module string) GormLogger {
	moduleLogger := l.log.With(zap.String("module", module))
	return GormLogger{
		ZapLogger:     moduleLogger,
		SlowThreshold: 200 * time.Millisecond,
	}
}

// Option 日志配置选项函数类型
type Option func(keyvals *[]any)

// WithModule 返回一个添加 module 字段的 Option
// module命名规范: "[组件]/[层]/[服务名]"
// 例如: "redis/data/micro-forge-service", "auth/biz/micro-forge-service"
func WithModule(module string) Option {
	return func(keyvals *[]any) {
		*keyvals = append(*keyvals, "module", module)
	}
}

// WithField 返回一个添加自定义字段的 Option
func WithField(key string, value any) Option {
	return func(keyvals *[]any) {
		*keyvals = append(*keyvals, key, value)
	}
}

// With 应用一系列 Option 到 logger 上
func With(logger log.Logger, opts ...Option) log.Logger {
	keyvals := make([]any, 0)
	for _, opt := range opts {
		opt(&keyvals)
	}
	return log.With(logger, keyvals...)
}
