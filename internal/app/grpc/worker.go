package grpc

import (
	xtremegrpc "github.com/globalxtreme/go-core/v2/grpc"
	"service/internal/app/grpc/server"
)

func Register(srv *xtremegrpc.GRPCServer) {
	srv.Register(
		&server.TestingServer{},
	)
}
