package server

import (
	"github.com/horonlee/krathub/internal/server/middleware"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(middleware.NewMiddlewareManager, NewRegistrar, NewGRPCServer, NewHTTPServer)
