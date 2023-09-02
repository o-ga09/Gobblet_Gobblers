package usecase

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestIsWin(t *testing.T) {
	t.Run("三目並べの勝敗を判定する(勝ち)",func(t *testing.T) {
		res:= IsWin()
		assert.Equal(t,res,true)
	})

	t.Run("三目並べの勝敗を判定する(負け)",func(t *testing.T) {
		res:= IsWin()
		assert.Equal(t,res,true)
	})

	t.Run("三目並べの勝敗を判定する(引き分け)",func(t *testing.T) {
		res:= IsWin()
		assert.Equal(t,res,true)
	})

	t.Run("縦一列が揃っているかを判定する",func(t *testing.T) {
		res:= CheckVertical()
		assert.Equal(t,res,true)
	})

	t.Run("横一列が揃っているかを判定する",func(t *testing.T) {
		res:= CheckHorizon()
		assert.Equal(t,res,true)
	})

	t.Run("斜め一列が揃っているかを判定する(右肩下がり)",func(t *testing.T) {
		res:= CheckCross()
		assert.Equal(t,res,true)
	})

	t.Run("斜め一列が揃っているかを判定する(右肩上がり)",func(t *testing.T) {
		res:= CheckCross()
		assert.Equal(t,res,true)
	})

	t.Run("コマを置けるか判定する",func(t *testing.T) {
		res:= IsEmpty()
		assert.Equal(t,res,true)
	})
}