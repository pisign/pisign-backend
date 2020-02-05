package game

import (
	"fmt"
)

// Player main interface for a Player object
type Player interface {
	fmt.Stringer
	makeMove() int
	chanMove() chan int
	chanGo() chan bool
	symbol() byte
}

type basePlayer struct {
	chMove     chan int
	chGo       chan bool
	name       string
	sym        byte
	difficulty PlayerType
	game       *Game
}

func (p *basePlayer) chanMove() chan int {
	return p.chMove
}

func (p *basePlayer) chanGo() chan bool {
	return p.chGo
}

func (p *basePlayer) String() string {
	return fmt.Sprintf("%s[%c]", p.name, p.sym)
}

func (p *basePlayer) init(name string, symbol byte, difficulty PlayerType, g *Game) {
	p.chMove = make(chan int, 1)
	p.chGo = make(chan bool, 1)
	p.name = name
	p.sym = symbol
	p.difficulty = difficulty
	p.game = g
}

func (p *basePlayer) symbol() byte {
	return p.sym
}

func play(p Player) {
	for <-p.chanGo() {
		move := makeMove(p)
		p.chanMove() <- move
	}
}

func makeMove(p Player) int {
	switch p := p.(type) {
	case *humanPlayer:
		return p.makeMove()
	case *cpuPlayer:
		return p.makeMove()
	default:
		panic(fmt.Sprintf("Unknown Player type %T!\n", p))
	}
}

func newPlayer(t PlayerType, name string, symbol byte, g *Game) Player {
	switch t {
	case Human:
		return newHumanPlayer(name, symbol, t, g)
	default:
		return newCPUPlayer(name, symbol, t, g)
	}
}
