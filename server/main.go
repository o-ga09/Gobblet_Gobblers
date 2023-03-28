package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"

	tictactoepb "main/pkg/grpc"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ROOM_ID_MAX = 100000
)

var count int

type TicTocToeGameServer struct {
	tictactoepb.UnimplementedGameServiceServer
	requests []*tictactoepb.GameRequest
	rooms []room
}

type room struct {
	room_id string
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
	count = 0
	for {
		in, err := stream.Recv()
		if err  == io.EOF {
			count++
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

			if count == 2 {
				err := stream.Send(&tictactoepb.GameResponse{
					PlayerName: res.GetPlayerName(),
					X: res.GetX(),
					Y: res.GetY(),
				})
				if err != nil {
					log.Fatal(err)
				}
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

func (s *TicTocToeGameServer) Greet(ctx context.Context,r *tictactoepb.GreetRequest) (*tictactoepb.GreetResponse,error) {
	log.Printf("Request from %s",r.GetMsg())
	return &tictactoepb.GreetResponse{Msg: fmt.Sprintf("Hello %s!",r.GetMsg())},nil
}

func (s *TicTocToeGameServer) AddRoom(ctx context.Context,_ *emptypb.Empty) (*tictactoepb.RoomInfo,error) {
	log.Printf("Add Room Request")

	
	roomId, err := rand.Int(rand.Reader,big.NewInt(ROOM_ID_MAX))
	if err !=  nil {
		return nil,err
	}

	s.rooms = append(s.rooms,room{room_id: roomId.String()})
	index, err := searchRoom(s.rooms,roomId.String())
	if err != nil {
		log.Fatal(err)
	}

	room := s.rooms[index]
	return &tictactoepb.RoomInfo{RoomId: room.room_id},nil
}

func (s *TicTocToeGameServer) GetRoomInfo(ctx context.Context,r *tictactoepb.RoomRequest) (*tictactoepb.RoomInfo,error) {
	log.Printf("Get Room request")

	index, err := searchRoom(s.rooms,r.RoomId)
	if err != nil {
		log.Fatal(err)
	}

	room := s.rooms[index]
	return &tictactoepb.RoomInfo{RoomId: room.room_id},nil
}

func (s *TicTocToeGameServer) GetRooms(ctx context.Context,_ *empty.Empty) (*tictactoepb.RoomList,error) {
	log.Printf("Get Rooms request")
	return &tictactoepb.RoomList{Room: buildRoomInfo(s.rooms)},nil
}

func buildRoomInfo(room []room) ([]*tictactoepb.RoomInfo){
	r := make([]*tictactoepb.RoomInfo,0)
	for _, v := range(room) {
		r = append(r,&tictactoepb.RoomInfo{RoomId: v.room_id})
	}
	return r
}

func searchRoom(r []room,room_id string) (int,error) {
	for i, v := range(r) {
		if v.room_id == room_id {
			return i,nil
		}
	}
	return -1,errors.New("not found")
}