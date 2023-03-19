package cmd

const (
	N = 3
	ROW_NUM = 3
	COLUMN_NUM = 3
	PLAYER1 = 1
	PALYER2 = 2
)


type Pos struct {
	attack int
	x int
	y int
}

// type Square struct {
// 	id int
// 	layer [3]int
// }

type Player interface {
	Init_Board(*[][]int)
	SetTurn(int)
	PrintBoard(*[][]int)
	Row_check(*[][]int) bool
	Column_check(*[][]int) bool
	Cross_check(*[][]int) bool
	Is_win(*[][]int) bool
	InputPlayer(*[][]int)
	Is_empty(*[][]int) bool
}