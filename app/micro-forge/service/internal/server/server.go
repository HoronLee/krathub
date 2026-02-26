package server

import (
	"github.com/google/wire"
	"github.com/horonlee/micro-forge/app/micro-forge/service/internal/server/middleware"
	"github.com/horonlee/micro-forge/pkg/governance/telemetry"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(middleware.ProviderSet, NewRegistrar, telemetry.NewMetrics, NewGRPCMiddleware, NewGRPCServer, NewHTTPMiddleware, NewHTTPServer)
