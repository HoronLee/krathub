package server

import (
	"github.com/google/wire"
	"github.com/horonlee/krathub/app/krathub/service/internal/server/middleware"
	"github.com/horonlee/krathub/pkg/governance/telemetry"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(middleware.ProviderSet, NewRegistrar, telemetry.NewMetrics, NewGRPCMiddleware, NewGRPCServer, NewHTTPMiddleware, NewHTTPServer)
