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