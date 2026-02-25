package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/horonlee/krathub/pkg/governance/telemetry"	
	"github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

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

	traceCleanup, err := telemetry.InitTracerProvider(bc.Trace, Name, bc.App.Env)
	if err != nil {
		panic(err)
	}
	defer traceCleanup()

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
