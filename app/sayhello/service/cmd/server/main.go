package main

import (
	"context"
	"flag"
	"os"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	_ "go.uber.org/automaxprocs"
)

var (
	Name     string
	Version  string
	flagconf string
	id, _    = os.Hostname()
	Metadata map[string]string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, reg registry.Registrar, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(Metadata),
		kratos.Logger(logger),
		kratos.Server(gs),
		kratos.Registrar(reg),
	)
}

func initTracerProvider(c *conf.Trace, env string) error {
	if c == nil || c.Endpoint == "" {
		return nil
	}

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(c.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithBatcher(exporter),
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

	bc, c, err := loadConfig()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	Name = bc.App.Name
	Version = bc.App.Version
	if Name == "" {
		Name = "sayhello.service"
	}
	if Version == "" {
		Version = "v0.1"
	}

	Metadata = bc.App.Metadata
	if Metadata == nil {
		Metadata = make(map[string]string)
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

	if err := initTracerProvider(bc.Trace, bc.App.Env); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Registry, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
