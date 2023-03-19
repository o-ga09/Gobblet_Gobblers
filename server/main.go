package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "main/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type sampleServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func main() {
	port := 8080
	l,err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("can not run server")
	}

	s := grpc.NewServer()

	hellopb.RegisterGreetingServiceServer(s,NewMyServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port : %v",port)
		s.Serve(l)
	}()

	quit := make(chan os.Signal,1)
	signal.Notify(quit,os.Interrupt)
	<- quit
	log.Printf("stopping grpc server")
	s.GracefulStop()
}

func (s *sampleServer) Hello(ctx context.Context,req *hellopb.HelloRequest) (*hellopb.HelloResponse,error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello %s!",req.GetName()),
	},nil
}

func NewMyServer() *sampleServer {
	return &sampleServer{}
}