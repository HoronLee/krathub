package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/horonlee/micro-forge/pkg/config"
	"github.com/horonlee/micro-forge/pkg/governance/telemetry"
	"github.com/horonlee/micro-forge/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	_ "go.uber.org/automaxprocs"
)

var (
	Name     = "sayhello.service"
	Version  = "v1.0.0"
	flagconf string
	id       string
	Metadata map[string]string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(l log.Logger, reg registry.Registrar, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(Metadata),
		kratos.Logger(l),
		kratos.Server(gs),
		kratos.Registrar(reg),
	)
}

func main() {
	flag.Parse()

	bc, c, err := config.LoadBootstrap(flagconf, Name)
	if err != nil {
		panic(err)
	}
	defer c.Close()

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

	appLogger := logger.NewLogger(&logger.Config{
		Env:        bc.App.Env,
		Level:      bc.App.Log.Level,
		Filename:   bc.App.Log.Filename,
		MaxSize:    bc.App.Log.MaxSize,
		MaxBackups: bc.App.Log.MaxBackups,
		MaxAge:     bc.App.Log.MaxAge,
		Compress:   bc.App.Log.Compress,
	})
	appLogger = log.With(
		appLogger,
		"service", Name,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	traceCleanup, err := telemetry.InitTracerProvider(bc.Trace, Name, bc.App.Env)
	if err != nil {
		panic(err)
	}
	defer traceCleanup()

	app, cleanup, err := wireApp(bc.Server, bc.Registry, bc.App, bc.Trace, bc.Metrics, appLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
