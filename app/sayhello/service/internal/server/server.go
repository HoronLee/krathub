package server

import (
	"github.com/google/wire"
	"github.com/horonlee/servora/pkg/governance/registry"
	"github.com/horonlee/servora/pkg/governance/telemetry"
)

var ProviderSet = wire.NewSet(registry.NewRegistrar, telemetry.NewMetrics, NewGRPCServer)
