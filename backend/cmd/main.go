package main

import (
	"fmt"
	"log"
	"main/api"
	"net"
	"os"
	"os/signal"

	hellopb "main/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 8080
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	hellopb.RegisterHellogRPCServiceServer(server,api.NewServer())

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