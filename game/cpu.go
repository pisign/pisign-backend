package game

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pisign/backend/easyrandom"
)

type cpuPlayer struct {
	basePlayer
	scanner *bufio.Scanner
}

func newCPUPlayer(name string, symbol byte, t PlayerType, g *Game) *cpuPlayer {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	cpu := &cpuPlayer{scanner: scanner}
	cpu.basePlayer.init(name, symbol, t, g)
	return cpu
}

func (p *cpuPlayer) String() string {
	return fmt.Sprintf("%s(CPU)", p.basePlayer.String())
}

func (p *cpuPlayer) makeMove() int {
	fmt.Printf("%s Thinking of move...\n", p)

	var free []int
	for i, cell := range p.game.board {
		if cell == 0 {
			free = append(free, i)
		}
	}
	move := free[easyrandom.RandomInt(0, int64(len(free)))]

	fmt.Printf("%s Made move %v\n", p, move)
	return move
}
