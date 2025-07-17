package main

import (
	"flag"
	"os"

	"krathub/internal/conf"
	"krathub/internal/server/middleware"
	"krathub/pkg"
	zapLog "krathub/pkg/log/zap"

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
	Name string = "krathub.service"
	// Version is the version of the compiled software.
	Version string = "v0.1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, reg registry.Registrar, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
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

	// 初始化一些外部包方法
	initComponents(bc)

	app, cleanup, err := wireApp(bc.Server, bc.Discovery, bc.Registry, bc.Data, bc.App, zapLog.Logger())
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initComponents(bc *conf.Bootstrap) {
	middleware.SetAppConf(bc.App)
	pkg.SetAppConf(bc.App)
}
