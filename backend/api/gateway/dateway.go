package gateway

import (
	"main/api/domain"
	"main/api/usecase/port"
)

type GatewayService struct {
	Port port.InputPort
}

func ProvideGatewayService(gateway port.InputPort) *GatewayService {
	return &GatewayService{Port: gateway}
}

func(g *GatewayService) Input(x,y,size,turn int) (int, error) {
	koma := domain.Koma{X: x,Y: y,Size: size,Turn: turn}
	res, err := g.Port.GameMain(koma)
	if err != nil {
		return -1, err
	}
	return res.Turn, nil
}