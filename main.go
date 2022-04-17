package main

import "fmt"

func main() {
	var board [][]int

	board = make([][]int, 3)
	Init_Board(&board)
	PrintBoard(&board)

}

func Init_Board(board *[][]int){
	fmt.Printf("%p\n",&board)
	for v := range (*board) {
		(*board)[v] = make([]int,3)
	}
	for i := 0;i < 3;i++ {
		for j := 0;j < 3;j++{
			(*board)[i][j] = 0
		}
	}
}

func PrintBoard(board *[][]int){
	for i := 0;i < 3;i++ {
		for j := 0;j < 3;j++{
			fmt.Print((*board)[i][j])
		}
		fmt.Printf("\n")
	}
}