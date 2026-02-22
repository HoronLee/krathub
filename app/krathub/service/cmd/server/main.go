package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"

	"github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "krathub.service"
	// Version is the version of the compiled software.
	Version = "v1.0.0"
	// flagconf is the config flag.
	flagconf string
	// id is the id of the instance.
	id string
	// Metadata is the service metadata.
	Metadata map[string]string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, reg registry.Registrar, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(Metadata),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
		kratos.Registrar(reg),
	)
}

// 设置全局trace
func initTracerProvider(c *conf.Trace, env string) error {
	if c == nil || c.Endpoint == "" {
		return nil
	}

	// 创建 exporter
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(c.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(exporter),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("exporter", "otlp"),
			attribute.String("env", env),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	flag.Parse()

	// 加载配置
	bc, c, err := loadConfig()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// 初始化服务名称、版本、元信息
	if bc.App.Name != "" {
		Name = bc.App.Name
	} else {
		bc.App.Name = Name
	}
	if bc.App.Version != "" {
		Version = bc.App.Version
	} else {
		bc.App.Version = Version
	}
	Metadata = bc.App.Metadata
	if Metadata == nil {
		Metadata = make(map[string]string)
	}

	hostname, _ := os.Hostname()
	id = fmt.Sprintf("%s-%s", Name, hostname)

	// 初始化日志
	// 如果未配置日志文件名，则使用默认值
	if bc.App.Log.Filename == "" {
		bc.App.Log.Filename = fmt.Sprintf("./logs/%s.log", Name)
	}
	log := logger.NewLogger(&logger.Config{
		Env:        bc.App.Env,
		Level:      bc.App.Log.Level,
		Filename:   bc.App.Log.Filename,
		MaxSize:    bc.App.Log.MaxSize,
		MaxBackups: bc.App.Log.MaxBackups,
		MaxAge:     bc.App.Log.MaxAge,
		Compress:   bc.App.Log.Compress,
	})

	// 初始化链路追踪
	if err := initTracerProvider(bc.Trace, bc.App.Env); err != nil {
		panic(err)
	}

	// 初始化服务
	app, cleanup, err := wireApp(bc.Server, bc.Discovery, bc.Registry, bc.Data, bc.App, bc.Trace, bc.Metrics, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 启动服务并且等待停止信号
	if err := app.Run(); err != nil {
		panic(err)
	}
}
