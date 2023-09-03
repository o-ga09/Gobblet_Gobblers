package main

import (
	"fmt"
	"main/api/domain"
	"main/api/usecase"
	"os"
	"strconv"
	"strings"
)

const (
	TURN1 = 1 //先攻
	TURN2 = 2 //後攻
)

func main() {
	service := usecase.GameService{}
	service.Init()

	var input_str string
	now_TURN := 1
	for {
		switch now_TURN {
		case TURN1:
			fmt.Println("Player 1 Input number koma position [x, y, size]")
			fmt.Fscan(os.Stdin,&input_str)
			args := strings.Split(input_str, ",")
			if len(args) != 3 {
				fmt.Println(args)
				os.Exit(1)
			}
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			size, _ := strconv.Atoi(args[2])
			koma := domain.Koma{X: x,Y: y,Size: size,Turn: now_TURN}
			if !service.IsEmpty(koma) {
				fmt.Println("Input another positon")
				continue
			}
			res, _ := service.Input(koma)
			now_TURN =res.Turn
		case TURN2:
			fmt.Println("Player 2 Input number koma position [x, y]")
			fmt.Fscan(os.Stdin,&input_str)
			args := strings.Split(input_str, ",")
			if len(args) != 3 {
				os.Exit(1)
			}
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			size, _ := strconv.Atoi(args[2])
			koma := domain.Koma{X: x,Y: y,Size: size,Turn: now_TURN}
			if !service.IsEmpty(koma) {
				fmt.Println("Input another positon")
				continue
			}
			res, _ := service.Input(koma)
			now_TURN =res.Turn
		}

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				fmt.Printf("%d",service.Board.BoardInfo[i][j].Turn)
			}
			fmt.Printf("\n")
		}

		if service.IsWin() {break}
	}
}