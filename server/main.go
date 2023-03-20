package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

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

func (s *sampleServer) HelloServerStream(req *hellopb.HelloRequest,stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5

	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello %s!",i,req.GetName()),
		});err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (s *sampleServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string,0)
	for {
		req, err := stream.Recv()
		if errors.Is(err,io.EOF) {
			message := fmt.Sprintf("Hello %v!",nameList)
			return stream.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			log.Fatal(err)
		}
		nameList = append(nameList,req.GetName())
	}
}

func (s *sampleServer) HelloBiStream(stream hellopb.GreetingService_HelloBiStreamServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err,io.EOF) {
			return nil
		}

		if err != nil {
			log.Fatal(err)
			return err
		}

		message := fmt.Sprintf("Hello %v!",req.GetName())
		if err := stream.Send(&hellopb.HelloResponse{
			Message: message,
		});err != nil {
			return err
		}
	}
}

func NewMyServer() *sampleServer {
	return &sampleServer{}
}