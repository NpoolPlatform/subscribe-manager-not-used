package api

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/subscribemgr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	npool.UnimplementedSubscribeManagerServer
}

func Register(server grpc.ServiceRegistrar) {
	npool.RegisterSubscribeManagerServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return npool.RegisterSubscribeManagerHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
