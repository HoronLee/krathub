package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/horonlee/krathub/pkg/governance/observability"
	"github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	_ "go.uber.org/automaxprocs"
)

var (
	Name     string
	Version  string
	flagconf string
	id       string
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
		Version = "v1.0.0"
	}

	Metadata = bc.App.Metadata
	if Metadata == nil {
		Metadata = make(map[string]string)
	}

	hostname, _ := os.Hostname()
	id = fmt.Sprintf("%s-%s", Name, hostname)

	log := logger.NewLogger(&logger.Config{
		Env:        bc.App.Env,
		Level:      bc.App.Log.Level,
		Filename:   bc.App.Log.Filename,
		MaxSize:    bc.App.Log.MaxSize,
		MaxBackups: bc.App.Log.MaxBackups,
		MaxAge:     bc.App.Log.MaxAge,
		Compress:   bc.App.Log.Compress,
	})

	traceCleanup, err := observability.InitTracerProvider(bc.Trace, Name, bc.App.Env)
	if err != nil {
		panic(err)
	}
	defer traceCleanup()

	app, cleanup, err := wireApp(bc.Server, bc.Registry, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
