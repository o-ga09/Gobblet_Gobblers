package main

import api "main/api/controller"

const (
	TURN1 = 1 //先攻
	TURN2 = 2 //後攻
)

func main() {
	api.Run()
}