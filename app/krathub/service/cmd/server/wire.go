//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
	"github.com/horonlee/krathub/app/krathub/service/internal/data"
	"github.com/horonlee/krathub/app/krathub/service/internal/server"
	"github.com/horonlee/krathub/app/krathub/service/internal/service"
	"github.com/horonlee/krathub/pkg/client"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Discovery, *conf.Registry, *conf.Data, *conf.App, *conf.Trace, *conf.Metrics, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, client.ProviderSet, newApp))
}
