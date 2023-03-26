package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	tictactoepb "main/pkg/grpc"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Client tictactoepb.GameServiceClient
	Conn *grpc.ClientConn
	Cancel chan struct{}
)

func Init_Board(board *[][]int){

	err := NewgRPCGameClient()
	if err != nil {
		log.Fatal(err)
	}

	for v := range (*board) {
		(*board)[v] = make([]int,COLUMN_NUM)
	}
	for i := 0;i < ROW_NUM;i++ {
		for j := 0;j < COLUMN_NUM;j++{
			(*board)[i][j] = 0
		}
	}
}

func (player *Pos) SetTurn(turn int) {
	player.attack = turn
}

func (player *Pos) PrintBoard(board *[][]int) {
	stream,err := Client.TicTacToeGame(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	res, err := stream.Recv()
	if err != nil {
		if errors.Is(err,io.EOF) {
			log.Fatal(err)
		}
	}

	pos,err := strconv.Atoi(res.GetPlayerName())
	if err != nil {
		log.Fatal(err)
	}
	(*board)[res.GetX()][res.GetY()] = pos
	for i := 0;i < ROW_NUM;i++ {
		for j := 0;j < COLUMN_NUM;j++{
			fmt.Print((*board)[i][j])
		}
		fmt.Printf("\n")
	}
}

func (player *Pos) Row_check(board *[][]int) bool {
	for i := 0;i < ROW_NUM;i++ {
		if (((*board)[i][0] == (*board)[i][1] && (*board)[i][1] == (*board)[i][2] && (*board)[i][0] == (*board)[i][2]) && (*board)[i][0] == player.attack){
			return true
		}
	}
	return false
}

func (player *Pos) Column_check(board *[][]int) bool{
	for i := 0;i < COLUMN_NUM;i++ {
		if (((*board)[0][i] == (*board)[1][i] && (*board)[1][i] == (*board)[2][i] && (*board)[0][i] == (*board)[2][i]) && (*board)[0][i] == player.attack){
			return true
		}
	}
	return false
}

func (player *Pos) Cross_check(board *[][]int) bool {
	if (((*board)[0][0] == (*board)[1][1] && (*board)[1][1] == (*board)[2][2] && (*board)[0][0] == (*board)[2][2]) && (*board)[0][0] == player.attack){
		return true
	}else if (((*board)[0][2] == (*board)[1][1] && (*board)[1][1] == (*board)[2][0] && (*board)[0][2] == (*board)[2][0]) && (*board)[0][2] == player.attack){
		return true
	}
	return false
}

func (player *Pos) Is_win(board *[][]int) bool {
	if (player.Row_check(board) || player.Column_check(board) || player.Cross_check(board)) {
		return true
	}
	return false
}

func (player *Pos) InputPlayer(board *[][]int) {

	stream, err := Client.TicTacToeGame(context.Background())
	if err != nil {
		log.Fatal(err)
	}



	fmt.Printf("Player %d (x y) : ",player.attack)
	fmt.Scanf("%d %d",&player.x,&player.y)
	fmt.Printf("x:%d,y:%d\n",player.x,player.y)
	if (player.Is_empty(board,player.x,player.y)) {
		if err := stream.Send(&tictactoepb.GameRequest{
			PlayerName: strconv.Itoa(player.attack),
			X: int32(player.x),
			Y: int32(player.y),
		});err != nil {
			log.Fatal(err)
		}

		if err := stream.CloseSend();err != nil {
			log.Fatal(err)
		}
	}
}

func (player *Pos) Is_empty(board *[][]int,x int,y int) bool {
	return (*board)[x][y] == 0
}

func NewgRPCGameClient() error {
	Cancel = make(chan struct{})
	//gRPC接続処理

	fmt.Println("start gRPC Client")

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	
	if err != nil {
		log.Printf("cannot open connection")
		return err
	}

	Client = tictactoepb.NewGameServiceClient(conn)
	return nil
	//gRPC接続処理ここまで
}