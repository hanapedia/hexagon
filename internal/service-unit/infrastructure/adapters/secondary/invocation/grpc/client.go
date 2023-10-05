package grpc

import (
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Connection *grpc.ClientConn
}

func NewGrpcClient(addr string) *GrpcClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// enable tracing
	if config.GetEnvs().TRACING {
		opts = append(opts, grpc.WithUnaryInterceptor(grpc.UnaryClientInterceptor(otelgrpc.UnaryClientInterceptor())))
		opts = append(opts, grpc.WithStreamInterceptor(grpc.StreamClientInterceptor(otelgrpc.StreamClientInterceptor())))
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Logger.Panicf("Failed to connect to Grpc server. err=%v, addr=%s", err, addr)
	}
	client := GrpcClient{
		Connection: conn,
	}

	// enable tracing

	return &client
}

func (gc *GrpcClient) Close() {
	gc.Connection.Close()
}
