package main

import (
	"context"
	"flag"
	"krathub/internal/conf"
	"krathub/pkg/logger"
	"os"
	"strings"

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
	Name string = "krathub.service"
	// Version is the version of the compiled software.
	Version string = "v0.1"
	// flagconf is the config flag.
	flagconf string
	// id is the id of the instance.
	id, _    = os.Hostname()
	Metadata map[string]string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
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
func initTracerProvider(c *conf.Trace) error {
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
			attribute.String("env", "dev"),
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

	// 初始化链路追踪
	if err := initTracerProvider(bc.Trace); err != nil {
		panic(err)
	}

	// 设置服务实例的元信息
	Metadata = bc.App.Metadata
	if Metadata == nil {
		Metadata = make(map[string]string)
	}
	// 从注册中心配置中提取 tags 信息（如果有）
	if bc.Registry != nil {
		switch r := bc.Registry.Registry.(type) {
		case *conf.Registry_Consul_:
			if len(r.Consul.Tags) > 0 {
				Metadata["tags"] = strings.Join(r.Consul.Tags, ";")
			}
			// 如果有其他注册中心（如 Nacos）需要处理类似 tags 的逻辑，可以在这里添加 case
		}
	}

	// 初始化日志
	log := logger.NewLogger(&logger.Config{
		Env:        bc.App.Env,
		Level:      bc.App.Log.Level,
		Filename:   bc.App.Log.Filename,
		MaxSize:    bc.App.Log.MaxSize,
		MaxBackups: bc.App.Log.MaxBackups,
		MaxAge:     bc.App.Log.MaxAge,
		Compress:   bc.App.Log.Compress,
	})

	// 初始化服务
	app, cleanup, err := wireApp(bc.Server, bc.Discovery, bc.Registry, bc.Data, bc.App, bc.Trace, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 启动服务并且等待停止信号
	if err := app.Run(); err != nil {
		panic(err)
	}
}
