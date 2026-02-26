package server

import (
	"github.com/google/wire"
	"github.com/horonlee/micro-forge/pkg/governance/telemetry"
)

var ProviderSet = wire.NewSet(NewRegistrar, telemetry.NewMetrics, NewGRPCServer)
