package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	tictactoepb "main/pkg/grpc"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Client tictactoepb.GameServiceClient
	Conn *grpc.ClientConn
	Cancel chan struct{}
	scanner bufio.Scanner
)

func Init(board *[][]int,player *Pos){

	err := NewgRPCGameClient()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Please enter name ->")
	scanner = *bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	stream, err := Client.TicTacToeGame(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&tictactoepb.GameRequest{
		PlayerName: name,
		X: 0,
		Y: 0,
		Attack: false,
	});err != nil {
		log.Fatal(err)
	}

	if err := stream.CloseSend();err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	var res *tictactoepb.GameResponse
	res = nil
	loading := []string{"m","a","t","c","h","i","n","g","・"}
	for i:=0;res == nil; {
		fmt.Printf("%s ",loading[i])
		time.Sleep(time.Second * 5)
		
		res, err = stream.Recv()
		if err != nil {
			if errors.Is(err,io.EOF) {
				log.Fatal(err)
			}
		}
		if i < len(loading) {i++}
	}
	fmt.Println()

	if res.GetPlayerName() == name {
		(*player).order = res.Attack
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
		if (((*board)[i][0] == (*board)[i][1] && (*board)[i][1] == (*board)[i][2] && (*board)[i][0] == (*board)[i][2]) && (*board)[i][0] == player.Attack){
			return true
		}
	}
	return false
}

func (player *Pos) Column_check(board *[][]int) bool{
	for i := 0;i < COLUMN_NUM;i++ {
		if (((*board)[0][i] == (*board)[1][i] && (*board)[1][i] == (*board)[2][i] && (*board)[0][i] == (*board)[2][i]) && (*board)[0][i] == player.Attack){
			return true
		}
	}
	return false
}

func (player *Pos) Cross_check(board *[][]int) bool {
	if (((*board)[0][0] == (*board)[1][1] && (*board)[1][1] == (*board)[2][2] && (*board)[0][0] == (*board)[2][2]) && (*board)[0][0] == player.Attack){
		return true
	}else if (((*board)[0][2] == (*board)[1][1] && (*board)[1][1] == (*board)[2][0] && (*board)[0][2] == (*board)[2][0]) && (*board)[0][2] == player.Attack){
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



	fmt.Printf("Player %d (x y) : ",player.Attack)
	fmt.Scanf("%d %d",&player.X,&player.Y)
	fmt.Printf("x:%d,y:%d\n",player.X,player.Y)
	if (player.Is_empty(board,player.X,player.Y)) {
		if err := stream.Send(&tictactoepb.GameRequest{
			PlayerName: strconv.Itoa(player.Attack),
			X: int32(player.X),
			Y: int32(player.Y),
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