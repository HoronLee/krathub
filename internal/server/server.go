package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegistrar, NewGRPCServer, NewHTTPServer)
