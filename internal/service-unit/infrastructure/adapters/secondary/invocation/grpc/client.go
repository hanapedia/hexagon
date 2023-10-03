package grpc

import (
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Connection *grpc.ClientConn
}

func NewGrpcClient(addr string) *GrpcClient {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(addr, opt)
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
