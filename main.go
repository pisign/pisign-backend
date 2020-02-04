package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type player interface {
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
	p.name = name
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

const (
	cellBlank byte = 0
	cellX     byte = 1
	cellO     byte = 2
)

type game struct {
	x     player
	o     player
	size  int
	board []byte
}

func newHumanVsHumanGame() *game {
	p1 := newHumanPlayer("P1[X]", 'X')
	p2 := newHumanPlayer("P2[O]", 'O')
	size := 3
	return &game{p1, p2, size, make([]byte, size*size)}
}

func (g *game) String() string {
	var str string
	for i, cell := range g.board {
		if cell == 0 {
			cell = '_'
		}
		str += string(cell)
		if (i+1)%g.size == 0 {
			str += "\n"
		}
	}
	return str
}

func (g *game) update(p player, move int) bool {
	g.board[move] = p.symbol()
	fmt.Println(g)
	return true
}

func (g *game) play(winner chan player) {
	activePlayer := g.x
	go g.x.play()
	go g.o.play()
	fmt.Println(g)
	for {
		fmt.Printf("Sending activation to %s...\n", activePlayer)
		activePlayer.chanGo() <- true
		fmt.Printf("Waiting for %s\n", activePlayer)
		move := <-activePlayer.chanMove()
		fmt.Printf("%v made move %v\n", activePlayer, move)
		g.update(activePlayer, move)
		switch activePlayer {
		case g.x:
			activePlayer = g.o
		default:
			activePlayer = g.x
		}
	}

}

func main() {
	g := newHumanVsHumanGame()
	winnerChan := make(chan player)
	g.play(winnerChan)
	winner := <-winnerChan
	fmt.Printf("Winner is %v!\n", winner)
}
