package client

import (
	"context"
	"fmt"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// client Client接口的实现
type client struct {
	dataCfg   *conf.Data
	traceCfg  *conf.Trace
	discovery registry.Discovery
	logger    log.Logger
}

// NewClient 创建客户端
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
		logger:    logger,
	}, nil
}

// CreateConn 创建指定类型的连接
func (c *client) CreateConn(ctx context.Context, connType ConnType, serviceName string) (Connection, error) {
	switch connType {
	case GRPC:
		return c.createGrpcConn(ctx, serviceName)
	default:
		return nil, fmt.Errorf("unsupported connection type: %s", connType)
	}
}

// createGrpcConn 创建gRPC连接
func (c *client) createGrpcConn(ctx context.Context, serviceName string) (Connection, error) {
	grpcConn, err := createGrpcConnection(ctx, serviceName, c.dataCfg, c.traceCfg, c.discovery, c.logger)
	if err != nil {
		return nil, err
	}

	return NewGrpcConn(grpcConn), nil
}
