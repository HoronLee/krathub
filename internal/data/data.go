package data

import (
	"github.com/horonlee/krathub/internal/client"
	"github.com/horonlee/krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewData, NewSayHelloRepo)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, clientFactory client.ClientFactory) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{
		log: log.NewHelper(logger),
	}, cleanup, nil
}
