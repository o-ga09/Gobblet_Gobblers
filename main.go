package main

import (
	"fmt"
	"main/cmd"
)

func main() {
	var player1 cmd.Pos
	var player2 cmd.Pos

	board := make([][]int, cmd.ROW_NUM)
	player1.Init_Board(&board)

	player1.SetTurn(cmd.PLAYER1)
	player2.SetTurn(cmd.PALYER2)
	for {
		player1.InputPlayer(&board)
		player1.PrintBoard(&board)
		if (player1.Is_win(&board)) {
			fmt.Printf("Player 1 is Win\n")
			break
		}
		player2.InputPlayer(&board)
		player2.PrintBoard(&board)
		if (player2.Is_win(&board)) {
			fmt.Printf("Player 2 is Win\n")
			break
		}
	}
}