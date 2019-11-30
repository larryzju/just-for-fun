package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
)

const escape = "\x1b"

const (
	BorderTopLeft     = '\u2554'
	BorderTopRight    = '\u2557'
	BorderBottomLeft  = '\u255a'
	BorderBottomRight = '\u255d'
	BorderHorization  = '\u2550'
	BorderVertical    = '\u2551'
)

const (
	Width  = 12
	Height = 12
)

type Block struct {
}

var defaultBlock = &Block{}

var (
	ShapeL = Shape{".X..X..XX", 3, defaultBlock}
	ShapeF = Shape{".XX.X..X.", 3, defaultBlock}
	ShapeT = Shape{"...XXX.X.", 3, defaultBlock}
	ShapeZ = Shape{"XX..XX...", 3, defaultBlock}
	Shape5 = Shape{".XX.X.XX.", 3, defaultBlock}
	ShapeO = Shape{"XXXX", 2, defaultBlock}
	ShapeI = Shape{"....XXXX.........", 4, defaultBlock}
)

type Shape struct {
	data  string
	n     int
	block *Block
}

func (s Shape) Rotate() Shape {
	bytes := make([]byte, len(s.data))
	for i := 0; i < s.n*s.n; i++ {
		r, c := i/s.n, i%s.n
		t := c*s.n + s.n - r - 1
		bytes[t] = s.data[i]
	}
	return Shape{string(bytes), s.n, s.block}
}

type Point struct {
	y, x int
}

func cleanScreen() {
	fmt.Printf("%s[s", escape)
}

func hideCursor() {
	// fmt.Printf("\e[?25l")
}

func moveCursor(line, col int) {
	fmt.Printf("%s[%d;%dH", escape, line, col)
}

type Game struct {
	Score  int
	Level  int
	Origin Point
	shape  Shape
	pos    Point
	board  []*Block
}

func (g *Game) drawBoard(topLeft Point, width, height int) {
	moveCursor(topLeft.x-1, topLeft.y-1)
	fmt.Printf("%c", BorderTopLeft)
	for i := 0; i < width; i++ {
		fmt.Printf("%c", BorderHorization)
	}
	fmt.Printf("%c", BorderTopRight)

	for i := 0; i < height; i++ {
		moveCursor(topLeft.y+i, topLeft.x-1)
		fmt.Printf("%c", BorderVertical)
		moveCursor(topLeft.y+i, topLeft.x+width)
		fmt.Printf("%c", BorderVertical)
	}

	moveCursor(topLeft.y+height, topLeft.x-1)
	fmt.Printf("%c", BorderBottomLeft)
	for i := 0; i < width; i++ {
		fmt.Printf("%c", BorderHorization)
	}
	fmt.Printf("%c", BorderBottomRight)
}

func (g *Game) flushScore() {
	moveCursor(g.Origin.y+1, g.Origin.x+Width+3)
	fmt.Printf("Score:%05d", g.Score)

	moveCursor(g.Origin.y+2, g.Origin.x+Width+3)
	fmt.Printf("Level:   %2d", g.Level)
}

func (g *Game) flushShape(o Point, s Shape, clear bool) {
	for i := 0; i < s.n*s.n; i++ {
		r, c := i/s.n, i%s.n
		if s.data[i] != '.' {
			moveCursor(g.Origin.y+o.y+r, g.Origin.x+o.x+c)
			char := s.data[i]
			if clear {
				char = ' '
			}
			fmt.Printf("%c", char)
		}
	}
}

func (g *Game) newShape(s Shape) {
	g.shape = s
	g.pos = Point{0, Width/2 - s.n/2}
	g.flushShape(g.pos, s, false)
	g.board = make([]*Block, Width*Height)
}

func (g *Game) Next(delta Point) {
	nextPos := Point{g.pos.y + delta.y, g.pos.x + delta.x}
	if g.isCollision(nextPos, g.shape) {
		return
	}
	g.flushShape(g.pos, g.shape, true)
	g.pos = nextPos
	g.flushShape(g.pos, g.shape, false)
}

func (g *Game) isCollision(pos Point, s Shape) bool {
	for i := 0; i < s.n*s.n; i++ {
		if s.data[i] == '.' {
			continue
		}

		r, c := i/s.n, i%s.n

		if pos.y+r < 0 || pos.y+r >= Height || pos.x+c < 0 || pos.x+c >= Width {
			return true
		}

		if g.board[(pos.y+r)*Width+pos.x+c] != nil {
			return true
		}
	}
	return false
}

type Key byte

const (
	Left Key = iota
	Right
	Down
	TurnLeft
	TurnRight
)

// ReadInput get input from keyboard
func ReadInput() <-chan Key {
	ch := make(chan Key, 3)
	go func() {
		for {
			input := make([]byte, 1)
			_, err := os.Stdin.Read(input)
			if err != nil {
				break
			}

			switch input[0] {
			case 'h':
				ch <- Left
			case 'l':
				ch <- Right
			case 'j':
				ch <- Down
			case 'k':
				ch <- TurnRight
			}
		}
	}()
	return ch
}

// Listen to the input and ticks
func (g *Game) Listen() {
	input := ReadInput()
	tick := time.Tick(1000 * time.Millisecond)
	for {
		select {
		case key := <-input:
			switch key {
			case Left:
				g.Next(Point{0, -1})
			case Right:
				g.Next(Point{0, 1})
			case Down:
				g.Next(Point{1, 0})
			case TurnRight:
				newShape := g.shape.Rotate()
				if !g.isCollision(g.pos, newShape) {
					g.flushShape(g.pos, g.shape, true)
					g.shape = newShape
					g.flushShape(g.pos, g.shape, false)
				}
			}
		case <-tick:
			if !g.isCollision(Point{g.pos.y + 1, g.pos.x}, g.shape) {
				g.Next(Point{1, 0})
			}
			g.flushScore()
		}
	}
}

func (g *Game) init() {
	mode, err := TCGetLocalMode(os.Stdin.Fd())
	if err != nil {
		log.Fatal(err)
	}

	mode &= ^uint64(syscall.ECHO | syscall.ICANON)
	log.Printf("%08X\n", mode)
	if err := TCSetLocalMode(os.Stdin.Fd(), mode); err != nil {
		log.Fatal(err)
	}

	cleanScreen()
	hideCursor()
	g.drawBoard(Point{10, 10}, Width, Height)
	g.flushScore()
	moveCursor(12, 12)
}

func main() {
	g := Game{Origin: Point{10, 10}}
	g.init()
	g.newShape(ShapeT)
	g.Listen()
	moveCursor(80, 0)
}
