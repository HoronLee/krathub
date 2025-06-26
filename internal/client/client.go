package client

import "github.com/google/wire"

// ProviderSet 是客户端工厂的依赖注入提供者集合
var ProviderSet = wire.NewSet(
	NewGrpcClientFactory,
)
