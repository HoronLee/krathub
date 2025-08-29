//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/horonlee/krathub/internal/biz"
	"github.com/horonlee/krathub/internal/client"
	"github.com/horonlee/krathub/internal/conf"
	"github.com/horonlee/krathub/internal/data"
	"github.com/horonlee/krathub/internal/server"
	"github.com/horonlee/krathub/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Discovery, *conf.Registry, *conf.Data, *conf.App, *conf.Trace, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, client.ProviderSet, newApp))
}
