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
)

const (
	ROOM_ID_MAX = 100000
)

type TicTocToeGameServer struct {
	tictactoepb.UnimplementedGameServiceServer
	rooms []room
}

type room struct {
	room_id string
	member [2]string
	contents []*gameInfo
}

type gameInfo struct {
	name string
	x int
	y int
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

	for {
		select {
		case v :=<- rcvCh:
			log.Printf("Received: [playername]%v, [X]%v, [Y]%v",v.GetPlayername(),v.GetX(),v.GetY())
			// チャットルームの探索
			index, err := searchRoom(s.rooms, v.GetRoomId())
			if err != nil {
				return err
			}
			// 対象チャットルーム
			targetRoom := s.rooms[index]
			msg, _ := latestMessage(targetRoom.contents)
			// クライアントへメッセージ送信
			if err := stream.Send(&tictactoepb.GameResponse{
				PlayerName: msg.name,
				X: int32(msg.x),
				Y: int32(msg.y),
				
			}); err != nil {
				return err
			}
			log.Printf("Sent : [playername]%v,[room id] %v, [X]%v, [Y]%v",v.GetPlayername(),v.GetRoomId(),v.GetX(),v.GetY())
		default:
			log.Printf("waiting")
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

		newR := &gameInfo{
			in.GetPlayername(),
			int(in.GetX()),
			int(in.GetY()),
		}

		index, err := searchRoom(s.rooms,in.GetRoomId())
		if err != nil {
			log.Fatal(err)
		}
		s.rooms[index].contents = append(s.rooms[index].contents, newR)
		ch <- &tictactoepb.GameRequest{
			RoomId: in.GetPlayername(),
			Playername: in.GetPlayername(),
			X: in.GetX(),
			Y: in.GetY(),
		}
	}
}

func (s *TicTocToeGameServer) Greet(ctx context.Context,r *tictactoepb.GreetRequest) (*tictactoepb.GreetResponse,error) {
	log.Printf("Request from %s",r.GetMsg())
	return &tictactoepb.GreetResponse{Msg: fmt.Sprintf("Hello %s!",r.GetMsg())},nil
}

func (s *TicTocToeGameServer) AddRoom(ctx context.Context,r *tictactoepb.RoomRequest) (*tictactoepb.RoomInfo,error) {
	log.Printf("Add Room Request")

	
	roomId, err := rand.Int(rand.Reader,big.NewInt(ROOM_ID_MAX))
	if err !=  nil {
		return nil,err
	}

	member := [2]string{}
	member[0] = r.GetPlayername()
	member[1] = ""
	s.rooms = append(s.rooms,room{room_id: roomId.String(),member: member})
	index, err := searchRoom(s.rooms,roomId.String())
	if err != nil {
		log.Fatal(err)
	}

	room := s.rooms[index]
	return &tictactoepb.RoomInfo{RoomId: room.room_id},nil
}

func (s *TicTocToeGameServer) JoinRoom(ctx context.Context,r *tictactoepb.RoomRequest) (*tictactoepb.RoomInfo,error) {
	log.Printf("Join Room %s Request",r.GetRoomId())

	index, err := searchRoom(s.rooms,r.GetRoomId())
	if err != nil {
		return nil, err
	}

	room := s.rooms[index]
	if room.member[1] != "" {
		return nil,errors.New("room is fulled")
	}else {
		room.member[1] = r.GetPlayername()
	}
	return &tictactoepb.RoomInfo{RoomId: room.room_id},nil
}

func (s *TicTocToeGameServer) GetRoomInfo(ctx context.Context,r *tictactoepb.RoomRequest) (*tictactoepb.RoomInfo,error) {
	log.Printf("Get Room request")

	index, err := searchRoom(s.rooms,r.RoomId)
	if err != nil {
		return nil,err
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

func latestMessage(message []*gameInfo) (gameInfo,error) {
	length := len(message)
	if length == 0 {
		return gameInfo{},errors.New("noot found")
	}
	return *message[length - 1],nil
}