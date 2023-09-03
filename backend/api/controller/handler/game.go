package handler

import (
	"context"
	"main/api/domain"
	"main/api/gateway"
	gamepb "main/pkg/grpc"
)
type GameServer struct {
	gateway gateway.GatewayService
	gamepb.UnimplementedGameServiceServer
}

func(s *GameServer) TicTacToeGame(ctx context.Context, req *gamepb.GameRequest) (*gamepb.GameResponse, error) {
	koma := domain.Koma{X: 1,Y: 1,Size: 1,Turn: 1}
	r, _ := s.gateway.Port.GameMain(koma)
	res := gamepb.GameResponse{X: int32(r.X),Y: int32(r.Y),Size: int32(r.Size),Turn: int32(r.Turn),IsFinished: r.IsFinished,IsDraw: r.IsDraw,NextTurn: int32(r.Next_Turn)}
	return &res, nil
}

func ProvideGameService(gateway gateway.GatewayService) *GameServer {
	return &GameServer{gateway: gateway}
}