package usecase

import (
	"main/api/domain"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestIsWin(t *testing.T) {
	service := GameService{}

	t.Run("ゲーム開始処理",func(t *testing.T) {
		service.Init()
		
		for _, rows := range service.Board.BoardInfo {
			for _, cols := range rows {
				assert.Equal(t,cols.Size,0)
				assert.Equal(t,cols.Turn,0)
			}
		}
		assert.Equal(t,len(service.Board.BoardInfo),3)
		assert.Equal(t,len(service.Board.BoardInfo[0]),3)
		assert.Equal(t,len(service.Board.BoardInfo[1]),3)
		assert.Equal(t,len(service.Board.BoardInfo[2]),3)
	})

	t.Run("三目並べの勝敗を判定する(勝ち)",func(t *testing.T) {
		res:= service.IsWin()
		assert.Equal(t,res,true)
	})

	t.Run("三目並べの勝敗を判定する(負け)",func(t *testing.T) {
		res:= service.IsWin()
		assert.Equal(t,res,true)
	})

	t.Run("三目並べの勝敗を判定する(引き分け)",func(t *testing.T) {
		res:= service.IsWin()
		assert.Equal(t,res,true)
	})

	t.Run("縦一列が揃っているかを判定する",func(t *testing.T) {
		service.Init()
		for i, rows := range service.Board.BoardInfo {
			for j := range rows {
				if j == 1 {
					service.Board.BoardInfo[i][j].Turn = 1
				}
			}
		}
		res:= service.CheckVertical()
		assert.Equal(t,res,true)
	})

	t.Run("横一列が揃っているかを判定する",func(t *testing.T) {
		service.Init()
		for i, rows := range service.Board.BoardInfo {
			for j := range rows {
				if i == 1 {
					service.Board.BoardInfo[i][j].Turn = 1
				}
			}
		}
		res:= service.CheckHorizon()
		assert.Equal(t,res,true)
	})

	t.Run("斜め一列が揃っているかを判定する(右肩下がり)",func(t *testing.T) {
		service.Init()
		for i, rows := range service.Board.BoardInfo {
			for j := range rows {
				if i == j {
					service.Board.BoardInfo[i][j].Turn = 1
				}
			}
		}
		res:= service.CheckCross()
		assert.Equal(t,res,true)
	})

	t.Run("斜め一列が揃っているかを判定する(右肩上がり)",func(t *testing.T) {
		service.Init()
		for i, rows := range service.Board.BoardInfo {
			for j := range rows {
				if 2 - i == j {
					service.Board.BoardInfo[i][j].Turn = 1
				}
			}
		}
		res:= service.CheckCross()
		assert.Equal(t,res,true)
	})

	t.Run("コマを置けるか判定する",func(t *testing.T) {
		service.Init()
		koma := domain.Koma{X: 1,Y: 1,Size: 2,Turn: 1}
		res:= service.IsEmpty(koma)
		assert.Equal(t,res,true)
	})

	t.Run("盤面にコマを入力する",func(t *testing.T) {
		service.Init()
		koma := domain.Koma{X: 0,Y: 0,Size: 1,Turn: 1}
		res, err := service.Input(koma)
		assert.Equal(t,res.Turn,2)
		assert.Equal(t,service.Board.BoardInfo[0][0].Size,1)
		assert.Equal(t,service.Board.BoardInfo[0][0].Turn,1)
		assert.Equal(t,err,nil)
	})
}