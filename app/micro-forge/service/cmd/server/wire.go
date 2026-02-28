//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/horonlee/micro-forge/api/gen/go/conf/v1"
	"github.com/horonlee/micro-forge/app/micro-forge/service/internal/biz"
	"github.com/horonlee/micro-forge/app/micro-forge/service/internal/data"
	"github.com/horonlee/micro-forge/app/micro-forge/service/internal/server"
	"github.com/horonlee/micro-forge/app/micro-forge/service/internal/service"
	"github.com/horonlee/micro-forge/pkg/bootstrap"
	"github.com/horonlee/micro-forge/pkg/transport/client"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Discovery, *conf.Registry, *conf.Data, *conf.App, *conf.Trace, *conf.Metrics, bootstrap.SvcIdentity, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, client.ProviderSet, newApp))
}
