package client

import (
	"context"
	"fmt"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	pkglogger "github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

type client struct {
	dataCfg   *conf.Data
	traceCfg  *conf.Trace
	discovery registry.Discovery
	logger    log.Logger
}

func NewClient(
	dataCfg *conf.Data,
	traceCfg *conf.Trace,
	discovery registry.Discovery,
	logger log.Logger,
) (Client, error) {
	return &client{
		dataCfg:   dataCfg,
		traceCfg:  traceCfg,
		discovery: discovery,
		logger:    pkglogger.WithModule(logger, "client/client/krathub-service"),
	}, nil
}

func (c *client) CreateConn(ctx context.Context, connType ConnType, serviceName string) (Connection, error) {
	switch connType {
	case GRPC:
		return c.createGrpcConn(ctx, serviceName)
	default:
		return nil, fmt.Errorf("unsupported connection type: %s", connType)
	}
}

func (c *client) createGrpcConn(ctx context.Context, serviceName string) (Connection, error) {
	grpcConn, err := createGrpcConnection(ctx, serviceName, c.dataCfg, c.traceCfg, c.discovery, c.logger)
	if err != nil {
		return nil, err
	}

	return NewGrpcConn(grpcConn), nil
}
