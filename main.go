package main

import "fmt"

const ROW_NUM = 3
const COLUMN_NUM = 3
const PLAYER1 = 1
const PALYER2 = 2

type pos struct {
	x int
	y int
}

func main() {
	var board [][]int

	board = make([][]int, ROW_NUM)
	Init_Board(&board)
	for {
		InputPlayer1(&board)
		PrintBoard(&board)
		if (Is_win(&board,PLAYER1)) {
			fmt.Printf("Player 1 is Win\n")
			break
		}
		InputPlayer2(&board)
		PrintBoard(&board)
		if (Is_win(&board,PALYER2)) {
			fmt.Printf("Player 2 is Win\n")
			break
		}
	}
}

func Init_Board(board *[][]int){
	fmt.Printf("%p\n",&board)
	for v := range (*board) {
		(*board)[v] = make([]int,COLUMN_NUM)
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

func Row_check(board *[][]int,player int) bool{
	for i := 0;i < ROW_NUM;i++ {
		if (((*board)[i][0] == (*board)[i][1] && (*board)[i][1] == (*board)[i][2] && (*board)[i][0] == (*board)[i][2]) && (*board)[i][0] == player){
			return true
		}
	}
	return false
}

func Column_check(board *[][]int,player int) bool{
	for i := 0;i < COLUMN_NUM;i++ {
		if (((*board)[0][i] == (*board)[1][i] && (*board)[1][i] == (*board)[2][i] && (*board)[0][i] == (*board)[2][i]) && (*board)[0][i] == player){
			return true
		}
	}
	return false
}

func Cross_check(board *[][]int,player int) bool{
	if (((*board)[0][0] == (*board)[1][1] && (*board)[1][1] == (*board)[2][2] && (*board)[0][0] == (*board)[2][2]) && (*board)[0][0] == player){
		return true
	}else if (((*board)[0][2] == (*board)[1][1] && (*board)[1][1] == (*board)[2][0] && (*board)[0][2] == (*board)[2][0]) && (*board)[0][2] == player){
		return true
	}
	return false
}

func Is_win(board *[][]int,player int) bool{
	if (Row_check(board,player) == true || Column_check(board,player) == true || Cross_check(board,player) == true) {
		return true
	}
	return false
}

func InputPlayer1(board *[][]int) {
	var P1 pos

	fmt.Printf("Player 1 (x y) : ")
	fmt.Scanf("%d %d",&P1.x,&P1.y)
	fmt.Printf("x:%d,y:%d\n",P1.x,P1.y)
	if (Is_empty(board,P1.x,P1.y)) {
		(*board)[P1.x][P1.y] = 1
	}
}

func InputPlayer2(board *[][]int) {
	var P2 pos

	fmt.Printf("Player 2 (x y) : ")
	fmt.Scanf("%d %d",&P2.x,&P2.y)
	if (Is_empty(board,P2.x,P2.y)) {
		(*board)[P2.x][P2.y] = 2
	}
}

func Is_empty(board *[][]int,x int,y int) bool {
	if ((*board)[x][y] != 0) {
		return false
	}
	return true
}