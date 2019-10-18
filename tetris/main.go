package main

import (
	"fmt"
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

var (
	ShapeL = Shape{".X..X..XX", 3}
	ShapeF = Shape{".XX.X..X.", 3}
	ShapeT = Shape{"...XXX.X.", 3}
	ShapeZ = Shape{"XX..XX...", 3}
	Shape5 = Shape{".XX.X.XX.", 3}
	ShapeO = Shape{"XXXX", 2}
	ShapeI = Shape{"....XXXX.........", 4}
)

type Shape struct {
	data string
	n    int
}

func (s Shape) Rotate() Shape {
	bytes := make([]byte, len(s.data))
	for i := 0; i < s.n*s.n; i++ {
		r, c := i/s.n, i%s.n
		t := c*s.n + s.n - r - 1
		bytes[t] = s.data[i]
	}
	return Shape{string(bytes), s.n}
}

type Point struct {
	y, x int
}

func cleanScreen() {
	fmt.Printf("%s[s", escape)
}

func moveCursor(line, col int) {
	fmt.Printf("%s[%d;%dH", escape, line, col)
}

type Block struct {
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
	g.flushShape(g.pos, g.shape, true)
	g.pos = Point{g.pos.y + delta.y, g.pos.x + delta.x}
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

func (g *Game) init() {
	cleanScreen()
	g.drawBoard(Point{10, 10}, Width, Height)
	g.flushScore()
	moveCursor(12, 12)
}

func main() {
	g := Game{Origin: Point{10, 10}}
	g.init()
	g.newShape(ShapeT)
	go func() {
		for {
			g.Score += 1
			time.Sleep(time.Millisecond * 500)
			if !g.isCollision(Point{g.pos.y + 1, g.pos.x}, g.shape) {
				g.Next(Point{1, 0})
			}
			g.flushScore()
		}
	}()

	time.Sleep(10 * time.Second)
	moveCursor(80, 0)
}
