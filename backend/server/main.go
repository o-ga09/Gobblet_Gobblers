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

	"main/config"
	tictactoepb "main/pkg/grpc"
	"main/store"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ROOM_ID_MAX = 100000
	EMPTY = 99999
)

type TicTocToeGameServer struct {
	tictactoepb.UnimplementedGameServiceServer
	room_id string
	rooms []room
	store store.KVS
}

type room struct {
	room_id string
	member [2]int32
	contents []gameInfo
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

	replyCh := make(chan *tictactoepb.GameResponse)
	go s.reply(replyCh,s.room_id,stream)

	for {
		select {
			case v := <- rcvCh:
				log.Printf("Received: [playername]%v,[room id]%v, [X]%v, [Y]%v",v.GetPlayername(),v.GetRoomId(),v.GetX(),v.GetY())
			case v := <- replyCh:
				log.Printf("Sent : [playername]%v [X]%v, [Y]%v",v.GetPlayerName(),v.GetX(),v.GetY())		
		}
	}
	return nil
}

func NewGameServer() *TicTocToeGameServer {

	config, err := config.New()
	if err != nil {
		log.Printf("failed to load config")
	}


	store, err := store.NewKVS(context.Background(),config)
	if err != nil {
		log.Printf("failed to connect redis")
	}

	return &TicTocToeGameServer{
		store: *store,
	}
}

func (s *TicTocToeGameServer)receive(ch chan<- *tictactoepb.GameRequest,stream tictactoepb.GameService_TicTacToeGameServer) {
	for {
		in, err := stream.Recv()
		if err  == io.EOF {
			continue
		}

		if err != nil {
			log.Fatalf("failed to Rcv(): %v",err)
		}

		newR := gameInfo{
			in.GetPlayername(),
			int(in.GetX()),
			int(in.GetY()),
		}

		index, err := searchRoom(s.rooms,in.GetRoomId())
		if err != nil {
			log.Fatalf("failed to look up data [room id] in receive : %v",err)
		}
		s.rooms[index].contents = append(s.rooms[index].contents, newR)
		ch <- &tictactoepb.GameRequest{
			RoomId: in.RoomId,
			Playername: in.GetPlayername(),
			X: in.GetX(),
			Y: in.GetY(),
		}
	}
}

func (s *TicTocToeGameServer) reply(ch chan <- *tictactoepb.GameResponse,room_id string, stream tictactoepb.GameService_TicTacToeGameServer) {
	// チャットルームの探索
	index, err := searchRoom(s.rooms, room_id)
	if err != nil {
		return
	}
	// 対象チャットルーム
	var current_count int
	previous_count := 0
	for {
		current_count = len(s.rooms[index].contents)
		if previous_count < current_count {
			if msg, err := latestMessage(s.rooms[index].contents);err != nil {
				continue
			} else {
				v :=&tictactoepb.GameResponse{
					PlayerName: msg.name,
					X: int32(msg.x),
					Y: int32(msg.y),
				}
				if err := stream.Send(v);err != nil {
					return
				}
				ch <- v
			}
		}
		previous_count = current_count
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

	member := [2]int32{}
	member[0] = int32(roomId.Int64())%2
	member[1] = EMPTY

	order := false
	if member[0] == 0 {order = true}
	
	s.room_id = roomId.String()
	s.rooms = append(s.rooms,room{room_id: roomId.String(),member: member})
	index, err := searchRoom(s.rooms,roomId.String())
	if err != nil {
		log.Fatalf("failed to look up date [room id] in add room : %v",err)
	}

	room := s.rooms[index]
	return &tictactoepb.RoomInfo{RoomId: room.room_id,Playername: room.member[0],Attack: order},nil
}

func (s *TicTocToeGameServer) JoinRoom(ctx context.Context,r *tictactoepb.RoomRequest) (*tictactoepb.RoomInfo,error) {
	log.Printf("Join Room %s Request",r.GetRoomId())

	index, err := searchRoom(s.rooms,r.GetRoomId())
	if err != nil {
		return nil, err
	}

	if s.rooms[index].member[1] != EMPTY {
		return nil,errors.New("room is fulled")
	}else {
		s.rooms[index].member[1] = 1 - s.rooms[index].member[0]
	}

	s.room_id = s.rooms[index].room_id
	order := false
	if s.rooms[index].member[1] == 0 {order = true}
	return &tictactoepb.RoomInfo{RoomId: s.rooms[index].room_id,Playername: s.rooms[index].member[1],Attack: order},nil
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

func latestMessage(message []gameInfo) (gameInfo,error) {
	length := len(message)
	if length == 0 {
		return gameInfo{},errors.New("not found")
	}
	return message[length - 1],nil
}