package server

import (
	"github.com/google/wire"
	"github.com/horonlee/servora/app/servora/service/internal/server/middleware"
	"github.com/horonlee/servora/pkg/governance/registry"
	"github.com/horonlee/servora/pkg/governance/telemetry"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(middleware.ProviderSet, registry.NewRegistrar, telemetry.NewMetrics, NewGRPCMiddleware, NewGRPCServer, NewHTTPMiddleware, NewHTTPServer)
