package grpc

import (
	"context"
	xtremegrpc "github.com/globalxtreme/go-core/v2/grpc"
	"service/internal/pkg/config"
	"service/internal/pkg/grpc/example"
	"time"
)

type TestingGRPC struct {
	xtremegrpc.GRPCClient
	Testing example.TestingServiceClient
}

func NewTestingGRPC(timeout ...time.Duration) (*TestingGRPC, context.CancelFunc) {
	client := TestingGRPC{}
	cleanup := client.RPCDialClient(config.TestingRPC, timeout...)

	client.Testing = example.NewTestingServiceClient(client.Conn)

	return &client, cleanup
}
