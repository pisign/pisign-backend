package main

import (
	"fmt"
	"time"

	"github.com/pisign/backend/game"
)

func main() {
	g := game.NewGame(game.Human, game.CPUEasy)
	winnerChan := make(chan game.Player)
	go g.Play(winnerChan)
	var winner game.Player
	for winner == nil {
		select {
		case winner = <-winnerChan:
			fmt.Printf("Winner is %v!\n", winner)
		default:
			//fmt.Println("Still waiting for winner...")
			time.Sleep(1 * time.Second)
		}
	}
}
