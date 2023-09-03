package api

import (
	"fmt"
	"log"
	"main/api/controller/handler"
	hellopb "main/pkg/grpc"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer() *handler.HelloServer {
	return &handler.HelloServer{}
}

func Run() {
	port := 8080
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	hellopb.RegisterHellogRPCServiceServer(server,NewServer())

	reflection.Register(server)

	go func() {
		log.Printf("server is started port %d\n", port)
		server.Serve(l)
	} ()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,os.Interrupt)
	<- quit
	log.Printf("stopping server ...\n")
	server.GracefulStop()
}