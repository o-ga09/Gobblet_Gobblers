package main

import (
	"fmt"
	"main/cmd"
)

func main() {
	defer cmd.Conn.Close()
	defer close(cmd.Cancel)
	//三目並べゲーム処理
	var player cmd.Pos

	board := make([][]int, cmd.ROW_NUM)
	cmd.Init(&board,&player)
	
	for {
		player.InputPlayer(&board)
		player.PrintBoard(&board)
		if (player.Is_win(&board)) {
			fmt.Printf("Player 1 is Win\n")
			break
		}
	}
	//三目並べゲーム処理ここまで
}