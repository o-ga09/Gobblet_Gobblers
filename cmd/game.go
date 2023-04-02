package cmd

var EMPTY = 9999

//Pos構造体を生成
func NewPlayer() (*Pos,error) {
	return &Pos{},nil
}

//盤面の初期化
func Init(board *[][]int){
	for v := range (*board) {
		(*board)[v] = make([]int,COLUMN_NUM)
	}
	for i := 0;i < ROW_NUM;i++ {
		for j := 0;j < COLUMN_NUM;j++{
			(*board)[i][j] = EMPTY
		}
	}
}

//勝利判定などゲームのコアな部分の処理はPosのメソッドとする
func (player *Pos) Row_check(board *[][]int) bool {
	for i := 0;i < ROW_NUM;i++ {
		if (((*board)[i][0] == (*board)[i][1] && (*board)[i][1] == (*board)[i][2] && (*board)[i][0] == (*board)[i][2]) && (*board)[i][0] != EMPTY){
			return true
		}
	}
	return false
}

func (player *Pos) Column_check(board *[][]int) bool{
	for i := 0;i < COLUMN_NUM;i++ {
		if (((*board)[0][i] == (*board)[1][i] && (*board)[1][i] == (*board)[2][i] && (*board)[0][i] == (*board)[2][i]) && (*board)[0][i] != EMPTY){
			return true
		}
	}
	return false
}

func (player *Pos) Cross_check(board *[][]int) bool {
	if (((*board)[0][0] == (*board)[1][1] && (*board)[1][1] == (*board)[2][2] && (*board)[0][0] == (*board)[2][2]) && (*board)[0][0] != EMPTY){
		return true
	}else if (((*board)[0][2] == (*board)[1][1] && (*board)[1][1] == (*board)[2][0] && (*board)[0][2] == (*board)[2][0]) && (*board)[0][2] != EMPTY){
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


func (player *Pos) Is_empty(board *[][]int,x int,y int) bool {
	return (*board)[x][y] == EMPTY
}