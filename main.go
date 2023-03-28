package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"main/client"
	"main/cmd"
	"os"
)

func main() {
	scanner := *bufio.NewScanner(os.Stdin)

	//gRPCクラインアント構造体を生成
	client, err := client.NewgRPCGameClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Conn.Close()

	//ルーム作成
	for client.RoomId == "" {
		fmt.Println("ルーム作成/他人のルームに参加を選択してください")
		fmt.Println("1:ルームを作成")
		fmt.Println("2:ルームに参加")
		scanner.Scan()
		in := scanner.Text()
	
		switch in {
		case "1":
			if err := client.CreateRoom(context.Background());err != nil {
				continue
			}
		case "2":
			roomId := scanner.Text()
			if err := client.JoinRoom(context.Background(),roomId);err != nil {
				continue
			}
		}
	}

	//ゲーム処理構造体を生成
	player, _ := cmd.NewPlayer()
	//三目並べ盤面生成
	board := make([][]int, cmd.ROW_NUM)
	//三目並べ盤面初期化
	cmd.Init(&board)
	
	//ゲーム処理メイン
	for {
		player.InputPlayer(&board)
		player.PrintBoard(&board)
		if (player.Is_win(&board)) {
			fmt.Printf("Player 1 is Win\n")
			break
		}
	}
}