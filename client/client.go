package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	tictactoepb "main/pkg/grpc"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GameClient struct {
	Client tictactoepb.GameServiceClient
	Conn *grpc.ClientConn
	Cancel chan struct{}
	Scanner bufio.Scanner
	Address string
	RoomId string
}

func (c *GameClient) GreetServer(ctx context.Context,name string) error {
	r, err := c.Client.Greetserver(ctx,&tictactoepb.GreetRequest{Msg: name})
	if err != nil {
		return err
	}
	log.Printf("%s",r.GetMsg())
	return nil
}

func (c *GameClient) CreateRoom(ctx context.Context) error {
	r, err := c.Client.AddRoom(ctx,&emptypb.Empty{})
	if err != nil {
		return err
	}
	c.RoomId = r.GetRoomId()
	log.Printf("create room >> %s",r.GetRoomId())
	return nil
}

func (c *GameClient) JoinRoom(ctx context.Context,id string) error {
	c.RoomId = id
	return nil
}

func (c *GameClient) GetRoom(ctx context.Context,id string) error {
	r, err := c.Client.GetRoomInfo(ctx,&tictactoepb.RoomRequest{
		RoomId: id,
	})
	if err != nil {
		return err
	}
	log.Printf("Room Information >> room id : %s",r.GetRoomId())
	return nil
}

func (c *GameClient) GetRooms(ctx context.Context) error {
	r, err := c.Client.GetRooms(ctx,&empty.Empty{})
	if err != nil {
		return err
	}
	for _, v := range (r.GetRoom()) {
		log.Printf("Name: %s",v.GetRoomId())
	}
	return nil
}

func NewgRPCGameClient() (*GameClient ,error){
	var err error
	client := &GameClient{}

	client.Cancel = make(chan struct{})

	fmt.Println("start gRPC Client")

	client.Address = "localhost:8080"
	client.Conn, err = grpc.Dial(
		client.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	
	if err != nil {
		log.Printf("cannot open connection")
		return nil,err
	}

	client.Client = tictactoepb.NewGameServiceClient(client.Conn)
	return client,nil
}