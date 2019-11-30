package main

import (
	"fmt"
	"time"
)

type Key byte

// Available keys
const (
	Left Key = iota
	Right
	Down
	TurnLeft
	TurnRight
)

// Input abstract the keyboard input devices
type Input interface {
	Init()
	Input() <-chan Key
}

// Display abstract the UI
type Display interface {
	Init()
	DrawTetris(t *Tetris, clear bool)
	UpdateScore(score, level int)
}

const escape = "\x1b"

const (
	BorderTopLeft     = '\u2554'
	BorderTopRight    = '\u2557'
	BorderBottomLeft  = '\u255a'
	BorderBottomRight = '\u255d'
	BorderHorization  = '\u2550'
	BorderVertical    = '\u2551'
)

type Block struct {
}

var defaultBlock = &Block{}

var shapes = []Shape{
	Shape{".X..X..XX", 3, defaultBlock},         // L
	Shape{".XXkX..X.", 3, defaultBlock},         // F
	Shape{"...XXX.X.", 3, defaultBlock},         // T
	Shape{"XX..XX...", 3, defaultBlock},         // Z
	Shape{".XX.X.XX.", 3, defaultBlock},         // 5
	Shape{"XXXX", 2, defaultBlock},              // o
	Shape{"....XXXX.........", 4, defaultBlock}, // I
}

// Shape of the tetris
type Shape struct {
	data  string
	n     int
	block *Block
}

// Rotate shape clockwise
func (s Shape) Rotate() Shape {
	bytes := make([]byte, len(s.data))
	for i := 0; i < s.n*s.n; i++ {
		r, c := i/s.n, i%s.n
		t := c*s.n + s.n - r - 1
		bytes[t] = s.data[i]
	}
	return Shape{string(bytes), s.n, s.block}
}

// Point for coordinator
type Point struct {
	y, x int
}

func moveCursor(line, col int) {
	fmt.Printf("%s[%d;%dH", escape, line, col)
}

// Tetris is a shape in specify position
type Tetris struct {
	shape Shape
	pos   Point
}

// Move tetris x blocks right and y blocks down
func (t *Tetris) Move(y, x int) *Tetris {
	return &Tetris{
		pos:   Point{t.pos.y + y, t.pos.x + x},
		shape: t.shape,
	}
}

// Rotate tetris clockwise
func (t *Tetris) Rotate() *Tetris {
	return &Tetris{
		pos:   t.pos,
		shape: t.shape.Rotate(),
	}
}

type Game struct {
	Input         Input
	Display       Display
	Width, Height int
	Score         int
	Level         int
	prev, cur     *Tetris
	board         []*Block
}

func (g *Game) FlushTetris() {
	if g.prev != nil {
		g.Display.DrawTetris(g.prev, true)
	}

	g.Display.DrawTetris(g.cur, false)
}

// NewTetris generate tetris at the top
func (g *Game) NewTetris(seed int) *Tetris {
	shape := shapes[seed%(len(shapes))]
	pos := Point{0, g.Width/2 - shape.n/2}
	return &Tetris{shape: shape, pos: pos}
}

func (g *Game) Transform(next *Tetris) {
	if g.isCollision(next) {
		return
	}

	g.cur, g.prev = next, g.cur
	g.FlushTetris()
}

func (g *Game) isCollision(t *Tetris) bool {
	pos, s := t.pos, t.shape
	for i := 0; i < s.n*s.n; i++ {
		if s.data[i] == '.' {
			continue
		}

		r, c := i/s.n, i%s.n

		if pos.y+r < 0 || pos.y+r >= g.Height || pos.x+c < 0 || pos.x+c >= g.Width {
			return true
		}

		if g.board[(pos.y+r)*g.Width+pos.x+c] != nil {
			return true
		}
	}
	return false
}

// Listen to the input and ticks
func (g *Game) Listen() {
	input := g.Input.Input()
	tick := time.Tick(1000 * time.Millisecond)
	for {
		select {
		case key := <-input:
			switch key {
			case Left:
				g.Transform(g.cur.Move(0, -1))
			case Right:
				g.Transform(g.cur.Move(0, 1))
			case Down:
				g.Transform(g.cur.Move(1, 0))
			case TurnRight:
				g.Transform(g.cur.Rotate())
			}
		case <-tick:
			g.Transform(g.cur.Move(1, 0))
		}
	}
}

func (g *Game) Init() {
	g.Input.Init()
	g.Display.Init()

	// init board
	g.board = make([]*Block, g.Width*g.Height)

	// init tetris
	// TODO random
	g.cur = g.NewTetris(0)
	g.FlushTetris()

	// show score board
	g.Display.UpdateScore(g.Score, g.Level)
}

func main() {
	width, height := 15, 15
	g := Game{
		Width:  width,
		Height: height,
		Input:  &TTYInput{},
		Display: &TTYDisplay{
			Origin: Point{10, 10},
			Width:  width,
			Height: height,
		},
	}

	g.Init()
	g.Listen()
}
