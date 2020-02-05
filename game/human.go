package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type humanPlayer struct {
	basePlayer
	scanner *bufio.Scanner
}

func newHumanPlayer(name string, symbol byte, t PlayerType, g *Game) *humanPlayer {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	human := &humanPlayer{scanner: scanner}
	human.basePlayer.init(name, symbol, t, g)
	return human
}

func (p *humanPlayer) String() string {
	return fmt.Sprintf("%s(Human)", p.basePlayer.String())
}

func (p *humanPlayer) makeMove() int {
	fmt.Printf("%s Please enter move: ", p)
	p.scanner.Scan()
	input := p.scanner.Text()
	move, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s made move %v\n", p, move)
	return move
}
