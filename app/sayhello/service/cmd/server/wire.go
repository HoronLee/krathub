//go:build wireinject
// +build wireinject

package main

import (
	"github.com/horonlee/servora/api/gen/go/conf/v1"
	"github.com/horonlee/servora/app/sayhello/service/internal/server"
	"github.com/horonlee/servora/app/sayhello/service/internal/service"
	"github.com/horonlee/servora/pkg/bootstrap"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Registry, *conf.App, *conf.Trace, *conf.Metrics, bootstrap.SvcIdentity, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, newApp))
}
