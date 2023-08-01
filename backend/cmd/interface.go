package cmd

const (
	N = 3
	ROW_NUM = 3
	COLUMN_NUM = 3
)


type Pos struct {
	Attack int
	X int
	Y int
	Order bool
}

// type Square struct {
// 	id int
// 	layer [3]int
// }

type Player interface {
	PrintBoard(*[][]int)
	Row_check(*[][]int) bool
	Column_check(*[][]int) bool
	Cross_check(*[][]int) bool
	Is_win(*[][]int) bool
	InputPlayer(*[][]int)
	Is_empty(*[][]int) bool
}