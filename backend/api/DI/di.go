package di

import (
	"main/api/controller/handler"
	"main/api/gateway"
	"main/api/usecase"
)

func DIcontainer() *handler.GameServer {
	usecase := usecase.ProvideGameService()
	gateway := gateway.ProvideGatewayService(usecase)
	handler := handler.ProvideGameService(*gateway)

	return handler
}