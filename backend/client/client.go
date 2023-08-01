package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"main/cmd"
	tictactoepb "main/pkg/grpc"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GameClient struct {
	Client tictactoepb.GameServiceClient
	Conn *grpc.ClientConn
	Cancel chan struct{}
	Scanner bufio.Scanner
	Address string
	RoomId string
	Player cmd.Pos
	Board [][]int
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
	r, err := c.Client.AddRoom(ctx,&tictactoepb.RoomRequest{RoomId: ""})
	if err != nil {
		return err
	}
	c.RoomId = r.GetRoomId()
	c.Player.Attack = int(r.GetPlayername())
	log.Printf("create room >> %s",r.GetRoomId())
	
	switch (c.Player.Attack) {
	case 0:
		c.Player.Order = true
		log.Printf("you are 先攻")
	case 1:
		c.Player.Order = false
		log.Printf("you are 後攻")
	}
	return nil
}

func (c *GameClient) JoinRoom(ctx context.Context,id string) error {
	r, err := c.Client.JoinRoom(ctx,&tictactoepb.RoomRequest{
		RoomId: id,
	})
	if err != nil {
		return err
	}
	c.RoomId = r.GetRoomId()
	c.Player.Attack = int(r.GetPlayername())
	log.Printf("joined room >> %s",r.GetRoomId())
	
	switch (c.Player.Attack) {
	case 0:
		c.Player.Order = true
		log.Printf("you are 先攻")
	case 1:
		c.Player.Order = false
		log.Printf("you are 後攻")
	}
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

func (c *GameClient) CommunicateWithServer() error {
	stream, err := c.Client.TicTacToeGame(context.Background())
	if err != nil {
		log.Println(err)
	}

	rcvCh := make(chan *tictactoepb.GameResponse)
	go c.ReceiveGameInfo(rcvCh,stream)

	ctx,cancel := context.WithCancel(context.TODO())
	defer cancel()

	if c.Player.Order {
		c.InputPlayer()
		if err := stream.Send(&tictactoepb.GameRequest{
			RoomId: c.RoomId,
			Playername: strconv.Itoa(c.Player.Attack),
			X: int32(c.Player.X),
			Y: int32(c.Player.Y),
		}); err != nil {
			log.Fatalf("failed to send : %v",err)
		}
	}

	for {
		select {
		case <- ctx.Done():
			fmt.Println("Done")
			stream.CloseSend()
			return nil
		case v := <-rcvCh:
			fmt.Printf("now , player %s are marked in [ %d, %d ]\n",(*v).GetPlayerName(),(*v).GetX(),(*v).GetY())
			
			if t, err := strconv.Atoi((*v).GetPlayerName());err != nil {
				continue
			} else {
				(c.Board)[(*v).GetX()][(*v).GetY()] = t
			}
			for i := 0;i < cmd.ROW_NUM;i++ {
				fmt.Printf("[ ")
				for j := 0;j < cmd.COLUMN_NUM;j++{
					if (c.Board)[i][j] == cmd.EMPTY {
						fmt.Printf(" - ")
					} else{fmt.Printf("%d",(c.Board)[i][j])}
				}
				fmt.Printf(" ]\n")
			}
			fmt.Printf("\n")
			
			if c.Player.Is_win(&c.Board) {
				fmt.Printf("Player %s is Win ! ",(*v).GetPlayerName())
				return nil
			}
			
			if strings.Compare(strconv.Itoa(c.Player.Attack),(*v).PlayerName)  == 0{
				c.Player.Order = false
				} else {
					c.Player.Order = true
					c.InputPlayer()
					if err := stream.Send(&tictactoepb.GameRequest{
						RoomId: c.RoomId,
					Playername: strconv.Itoa(c.Player.Attack),
					X: int32(c.Player.X),
					Y: int32(c.Player.Y),
				}); err != nil {
					log.Fatalf("failed to send : %v",err)
				}
			}
		}
	}
}

func (c *GameClient) InputPlayer() {
	//プレイヤーの入力の受付
	for {
		switch (c.Player.Attack) {
		case 0:
			fmt.Printf("先攻の入力 x y : ")
		case 1:
			fmt.Printf("後攻の入力 x y : ")
		}
		fmt.Scanf("%d %d",&c.Player.X,&c.Player.Y)
		fmt.Printf("x:%d,y:%d\n",c.Player.X,c.Player.Y)
		if c.Player.Is_empty(&c.Board,c.Player.X,c.Player.Y) {
			(c.Board)[c.Player.X][c.Player.Y] = c.Player.Attack
			return
		}
		fmt.Println("すでに置かれた場所です")
	}
}

func (c *GameClient) ReceiveGameInfo(ch chan<- *tictactoepb.GameResponse,stream tictactoepb.GameService_TicTacToeGameClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			close(ch)
			return
		}

		if err != nil {
			log.Fatalf("failed to receive : %v",err)
		}
		ch<- in
	}
}

func NewgRPCGameClient() (*GameClient ,error){
	var err error
	client := &GameClient{}

	//三目並べ盤面生成
	client.Board = make([][]int, cmd.ROW_NUM)
	//三目並べ盤面初期化
	cmd.Init(&client.Board)

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