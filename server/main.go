package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	tictactoepb "main/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type TicTocToeGameServer struct {
	tictactoepb.UnimplementedGameServiceServer
	requests []*tictactoepb.GameRequest
}

func main() {
	port := 8080
	l,err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("can not run server")
	}

	s := grpc.NewServer()

	tictactoepb.RegisterGameServiceServer(s,NewGameServer())

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

func (s *TicTocToeGameServer) TicTacToeGame(stream tictactoepb.GameService_TicTacToeGameServer) error {
	rcvCh := make(chan *tictactoepb.GameRequest)
	go s.receive(rcvCh,stream)

	replyCh := make(chan *tictactoepb.GameRequest)
	go s.reply(replyCh,stream)

	for {
		select {
		case v :=<- rcvCh:
			log.Printf("Received: [playername]%v, [X]%v, [Y]%v",v.GetPlayerName(),v.GetX(),v.GetY())
		case v :=<- replyCh:
			log.Printf("Sent : [playername]%v, [X]%v, [Y]%v",v.GetPlayerName(),v.GetX(),v.GetY())
		}
	}
}

func NewGameServer() *TicTocToeGameServer {
	return &TicTocToeGameServer{}
}

func (s *TicTocToeGameServer)receive(ch chan<- *tictactoepb.GameRequest,stream tictactoepb.GameService_TicTacToeGameServer) {
	for {
		in, err := stream.Recv()
		if err  == io.EOF {
			continue
		}
		newR := &tictactoepb.GameRequest{
			PlayerName: in.GetPlayerName(),
			X: in.GetX(),
			Y: in.GetY(),
		}
		s.requests = append(s.requests, newR)
		ch <- newR
	}
}

func (s *TicTocToeGameServer)reply(ch chan<- *tictactoepb.GameRequest,stream tictactoepb.GameService_TicTacToeGameServer) {

	for {
		if s.requests != nil && len(s.requests) != 0 {
			res := s.requests[0]

			if err := stream.Send(&tictactoepb.GameResponse{
				PlayerName: res.GetPlayerName(),
				X: res.GetX(),
				Y: res.GetY(),
			});err != nil {
				log.Fatal(err)
			}

			if len(s.requests) >= 2 {
				s.requests = s.requests[1:]
				ch <- res
			} else {
				s.requests = nil
				ch <- res
			}
		}
	}
}