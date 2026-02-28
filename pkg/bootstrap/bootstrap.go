package bootstrap

import (
	"fmt"
	"os"

	conf "github.com/horonlee/micro-forge/api/gen/go/conf/v1"
	"github.com/horonlee/micro-forge/pkg/config"
	"github.com/horonlee/micro-forge/pkg/governance/telemetry"
	"github.com/horonlee/micro-forge/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

type SvcIdentity struct {
	Name     string
	Version  string
	ID       string
	Metadata map[string]string
}

type Runtime struct {
	Bootstrap *conf.Bootstrap
	Identity  SvcIdentity
	Logger    log.Logger

	configCloser func()
	traceCloser  func()
}

func NewRuntime(configPath, name, version string) (*Runtime, error) {
	bc, c, err := config.LoadBootstrap(configPath, name)
	if err != nil {
		return nil, err
	}

	if bc.App == nil {
		bc.App = &conf.App{}
	}

	hostname, _ := os.Hostname()
	identity := resolveServiceIdentity(name, version, hostname, bc.App)

	appLog := bc.App.GetLog()

	appLogger := logger.NewLogger(&logger.Config{
		Env:        bc.App.Env,
		Level:      appLog.GetLevel(),
		Filename:   appLog.GetFilename(),
		MaxSize:    appLog.GetMaxSize(),
		MaxBackups: appLog.GetMaxBackups(),
		MaxAge:     appLog.GetMaxAge(),
		Compress:   appLog.GetCompress(),
	})
	appLogger = log.With(
		appLogger,
		"service", identity.Name,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	traceCleanup, err := telemetry.InitTracerProvider(bc.Trace, identity.Name, bc.App.Env)
	if err != nil {
		c.Close()
		return nil, err
	}

	return &Runtime{
		Bootstrap: bc,
		Identity:  identity,
		Logger:    appLogger,
		configCloser: func() {
			_ = c.Close()
		},
		traceCloser: traceCleanup,
	}, nil
}

func (r *Runtime) Close() {
	if r == nil {
		return
	}
	if r.traceCloser != nil {
		r.traceCloser()
	}
	if r.configCloser != nil {
		r.configCloser()
	}
}

func Run(app *kratos.App) error {
	return app.Run()
}

func resolveServiceIdentity(defaultName, defaultVersion, hostname string, app *conf.App) SvcIdentity {
	name := defaultName
	version := defaultVersion
	metadata := make(map[string]string)

	if app != nil {
		if app.Name != "" {
			name = app.Name
		} else {
			app.Name = name
		}
		if app.Version != "" {
			version = app.Version
		} else {
			app.Version = version
		}
		if app.Metadata == nil {
			app.Metadata = metadata
		} else {
			metadata = app.Metadata
		}
	}

	id := fmt.Sprintf("%s-%s", name, hostname)
	return SvcIdentity{
		Name:     name,
		Version:  version,
		ID:       id,
		Metadata: metadata,
	}
}
