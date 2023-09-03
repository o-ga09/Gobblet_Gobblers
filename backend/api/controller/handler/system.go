package handler

import (
	"context"
	hellopb "main/pkg/grpc"
)

type HelloServer struct {
	hellopb.UnimplementedHellogRPCServiceServer
}

func(s *HelloServer) Greetserver(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	res := hellopb.HelloResponse{Msg: "Hello " + req.Msg + " !"}
	return &res, nil
}

func ProvideHealthCheckService() *HelloServer {
	return &HelloServer{}
}