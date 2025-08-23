//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"krathub/internal/biz"
	"krathub/internal/client"
	"krathub/internal/conf"
	"krathub/internal/data"
	"krathub/internal/server"
	"krathub/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Discovery, *conf.Registry, *conf.Data, *conf.App, *conf.Trace, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, client.ProviderSet, newApp))
}
