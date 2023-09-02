package domain

type Board struct {
	BoardInfo [][]BoardInfo
}

type Koma struct {
	Turn int
	X int
	Y int
	Size int
}

type BoardInfo struct {
	Size int
	Turn int
}

func (b *Board) Init() {
	for _, rows := range b.BoardInfo {
		for _, cols := range rows {
			cols.Size = 0
			cols.Turn = 0
		}
	}
}