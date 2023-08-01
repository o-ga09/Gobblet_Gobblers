package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"main/client"
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
			scanner.Scan()
			roomId := scanner.Text()
			if err := client.JoinRoom(context.Background(),roomId);err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	
	//ゲーム処理メイン
	if err := client.CommunicateWithServer();err != nil {
		fmt.Printf("err : %v",err)
	}
}