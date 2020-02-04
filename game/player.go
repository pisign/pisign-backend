package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Player main interface for a Player object
type Player interface {
	fmt.Stringer
	play() bool
	makeMove() int
	chanMove() chan int
	chanGo() chan bool
	symbol() byte
}

type basePlayer struct {
	chMove chan int
	chGo   chan bool
	name   string
	sym    byte
}

func (p *basePlayer) chanMove() chan int {
	return p.chMove
}

func (p *basePlayer) chanGo() chan bool {
	return p.chGo
}

func (p *basePlayer) String() string {
	return p.name
}

func (p *basePlayer) init(name string, symbol byte) {
	p.chMove = make(chan int, 1)
	p.chGo = make(chan bool, 1)
	p.name = fmt.Sprintf("%s[%c]", name, symbol)
	p.sym = symbol
}

func (p *basePlayer) symbol() byte {
	return p.sym
}

type humanPlayer struct {
	basePlayer
	scanner *bufio.Scanner
}

func newHumanPlayer(name string, symbol byte) *humanPlayer {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	human := &humanPlayer{scanner: scanner}
	human.basePlayer.init(name, symbol)
	return human
}

func (p *humanPlayer) play() bool {
	for {
		<-p.chanGo()
		move := p.makeMove()
		p.chanMove() <- move
	}
}

func (p *humanPlayer) makeMove() int {
	fmt.Printf("Make move: ")
	p.scanner.Scan()
	input := p.scanner.Text()
	move, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Made move %v\n", move)
	return move
}
