package usecase

import (
	"fmt"
	"main/api/domain"
)

type GameService struct {
	Board *domain.Board
}

const (
	ROW_NUM = 3
	COLUMN_NUM= 3
	EMPTY = 0
)

func ProvideGameService() *GameService {
	return &GameService{}
}

func(g *GameService) GameMain(koma domain.Koma) (*domain.GameResponse, error) {
	res := domain.GameResponse{
		X: 1,
		Y: 1,
		Turn: 1,
		Size: 1,
		IsFinished: false,
		IsDraw: false,
		Next_Turn: 2,
	}
	return &res, nil
}

func(g *GameService) Init() {
	board := make([][]domain.BoardInfo,3)
	for i := 0; i < 3; i++ {
		board[i] = make([]domain.BoardInfo,3)
	}
	g.Board = &domain.Board{BoardInfo: board}
	g.Board.Init()
}

func(g *GameService) IsWin() bool {
	if g.CheckVertical() || g.CheckHorizon() || g.CheckCross() {
		return true
	}

	return false
}

func(g *GameService) IsEmpty(koma domain.Koma) bool {
	return g.Board.BoardInfo[koma.X][koma.Y].Turn != koma.Turn && g.Board.BoardInfo[koma.X][koma.Y].Size < koma.Size
}

func(g *GameService) CheckCross() bool {
	if ((g.Board.BoardInfo[0][0].Turn == g.Board.BoardInfo[1][1].Turn && g.Board.BoardInfo[1][1].Turn == g.Board.BoardInfo[2][2].Turn && g.Board.BoardInfo[0][0].Turn == g.Board.BoardInfo[2][2].Turn) && g.Board.BoardInfo[0][0].Turn != EMPTY){
		return true
	}else if ((g.Board.BoardInfo[0][2].Turn == g.Board.BoardInfo[1][1].Turn && g.Board.BoardInfo[1][1].Turn == g.Board.BoardInfo[2][0].Turn && g.Board.BoardInfo[0][2].Turn == g.Board.BoardInfo[2][0].Turn) && g.Board.BoardInfo[0][2].Turn != EMPTY){
		return true
	}

	return false
}

func(g *GameService) CheckHorizon() bool {
	for i := 0;i < ROW_NUM;i++ {
		if ((g.Board.BoardInfo[i][0].Turn == g.Board.BoardInfo[i][1].Turn && g.Board.BoardInfo[i][1].Turn == g.Board.BoardInfo[i][2].Turn && g.Board.BoardInfo[i][0].Turn == g.Board.BoardInfo[i][2].Turn) && g.Board.BoardInfo[i][0].Turn != EMPTY){
			return true
		}
	}

	return false
}

func(g *GameService) CheckVertical() bool {
	for i := 0;i < COLUMN_NUM;i++ {
		if (((g.Board.BoardInfo[0][i].Turn) == g.Board.BoardInfo[1][i].Turn && g.Board.BoardInfo[1][i].Turn == g.Board.BoardInfo[2][i].Turn && g.Board.BoardInfo[0][i].Turn == g.Board.BoardInfo[2][i].Turn) && g.Board.BoardInfo[0][i].Turn != EMPTY){
			return true
		}
	}
	
	return false
}

func(g *GameService) Input(koma domain.Koma) (domain.Koma, error) {
	next_turn := 0
	if koma.X >= ROW_NUM || koma.Y >= COLUMN_NUM {
		return domain.Koma{}, fmt.Errorf("out range index: %d,%d",koma.X,koma.Y)
	}

	if !g.IsEmpty(koma) {
		return  domain.Koma{}, fmt.Errorf("DRAW")
	}

	g.Board.BoardInfo[koma.X][koma.Y].Size = koma.Size
	g.Board.BoardInfo[koma.X][koma.Y].Turn = koma.Turn
	if koma.Turn == 1 {
		next_turn = 2
	} else if koma.Turn == 2 {
		next_turn = 1
	} 
	return domain.Koma{X: koma.X, Y: koma.Y, Size: koma.Size, Turn: next_turn}, nil
}