package handler

import (
	"context"
	"main/api/usecase/port"
	gamepb "main/pkg/grpc"
)
type GameServer struct {
	port.InputPort
	gamepb.UnimplementedGameServiceServer
}

func(s *GameServer) TicTacToeGame(ctx context.Context, req *gamepb.GameRequest) (*gamepb.GameResponse, error) {
	res := gamepb.GameResponse{
		X: 1,
		Y: 2,
		Size: 1,
		Turn: 1,
	}
	return &res, nil
}

func ProvideGameService() *GameServer {
	return &GameServer{}
}