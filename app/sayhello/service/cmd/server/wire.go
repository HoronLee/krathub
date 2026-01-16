//go:build wireinject
// +build wireinject

package main

import (
	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/app/sayhello/service/internal/server"
	"github.com/horonlee/krathub/app/sayhello/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, newApp))
}
