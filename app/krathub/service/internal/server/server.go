package server

import (
	"github.com/horonlee/krathub/app/krathub/service/internal/server/middleware"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(middleware.ProviderSet, NewRegistrar, NewGRPCMiddleware, NewGRPCServer, NewHTTPMiddleware, NewHTTPServer, NewMetrics)
