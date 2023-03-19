package cmd

import "fmt"

func (player *Pos) Init_Board(board *[][]int){
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

	fmt.Printf("Player %d (x y) : ",player.attack)
	fmt.Scanf("%d %d",&player.x,&player.y)
	fmt.Printf("x:%d,y:%d\n",player.x,player.y)
	if (player.Is_empty(board,player.x,player.y)) {
		(*board)[player.x][player.y] = player.attack
	}
}

func (player *Pos) Is_empty(board *[][]int,x int,y int) bool {
	return (*board)[x][y] == 0
}