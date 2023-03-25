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

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type sampleServer struct {
	hellopb.UnimplementedGreetingServiceServer
	requests []*hellopb.MessageRequest
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

func (s *sampleServer)CreateMessage(ctx context.Context,r *hellopb.MessageRequest) (*hellopb.MessageResponse,error) {
	log.Printf("Received %v",r.GetMessage())
	newR := &hellopb.MessageRequest{
		Name: r.GetName(),
		Message: r.GetMessage(),
		CreatedAt: r.GetCreatedAt(),
	}
	s.requests = append(s.requests,newR)
	return &hellopb.MessageResponse{Message:r.GetMessage()},nil
}

func (s *sampleServer) GetMessages(_ *empty.Empty,stream hellopb.GreetingService_GetMessagesServer) error {
	for _, r := range s.requests {
		if err := stream.Send(&hellopb.MessageResponse{
			Message: r.GetMessage(),
		});err != nil {
			return err
		}
	}

	previousCount := len(s.requests)

	for {
		currentCount := len(s.requests)
		if previousCount < currentCount {
			r := s.requests[currentCount - 1]
			log.Printf("Sent: %v",r.GetMessage())
			if err := stream.Send(&hellopb.MessageResponse{
				Name: r.GetName(),
				Message: r.GetMessage(),
				CreatedAt: r.GetCreatedAt(),
			});err != nil {
				return err
			}
		}
		previousCount = currentCount
	}
}

func (s *sampleServer) Chat(stream hellopb.GreetingService_ChatServer) error {
	rcvCh := make(chan *hellopb.MessageRequest)
	go s.receive(rcvCh,stream)

	replyCh := make(chan *hellopb.MessageRequest)
	go s.reply(replyCh,stream)

	for {
		select {
		case v :=<- rcvCh:
			log.Printf("Received: [message]%v, [user]%v",v.GetMessage(),v.GetName())
		case v :=<- replyCh:
			log.Printf("Sent : [message]%v, [user]%v",v.GetMessage(),v.GetName())
			if err := stream.Send(&hellopb.MessageResponse{
				Name: v.GetName(),
				Message: v.GetMessage(),
				CreatedAt: v.GetCreatedAt(),
			});err != nil {
				return err
			}

		}
	}
}

func NewMyServer() *sampleServer {
	return &sampleServer{}
}

func (s *sampleServer)receive(ch chan<- *hellopb.MessageRequest,stream hellopb.GreetingService_ChatServer) {
	for {
		in, err := stream.Recv()
		if err  == io.EOF {
			continue
		}
		newR := &hellopb.MessageRequest{
			Name: in.GetName(),
			Message: in.GetMessage(),
			CreatedAt: in.GetCreatedAt(),
		}
		s.requests = append(s.requests,newR)
		ch <- newR
	}
}

func (s *sampleServer)reply(ch chan<- *hellopb.MessageRequest,stream hellopb.GreetingService_ChatServer) {
	for _, r := range s.requests {
		if err := stream.Send(&hellopb.MessageResponse{
			Name: r.GetName(),
			Message: r.GetMessage(),
			CreatedAt: r.GetCreatedAt(),
		});err != nil {
			log.Fatal(err)
		}
	}

	previousCount := len(s.requests)

	for {
		currentCount := len(s.requests)
		if previousCount < currentCount {
			r := s.requests[currentCount - 1]
			ch <- r
		}
		previousCount = currentCount
	}
}