package main

import (
	"fmt"
	"log"
	"math/rand"
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
	DrawBlocks(blocks []*Block)
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
	Shape{".XX.X..X.", 3, defaultBlock},         // F
	Shape{"...XXX.X.", 3, defaultBlock},         // T
	Shape{"XX..XX...", 3, defaultBlock},         // Z
	Shape{".XXXX....", 3, defaultBlock},         // 5
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

type Ticker struct {
	*time.Ticker
	Duration time.Duration
}

func NewTicker(duration time.Duration) *Ticker {
	return &Ticker{
		Duration: duration,
		Ticker:   time.NewTicker(duration),
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

	tickerDuration time.Duration

	over   chan struct{}
	score  chan int
	ticker *Ticker
}

func (g *Game) FlushTetris() {
	if g.prev != nil {
		g.Display.DrawTetris(g.prev, true)
	}

	g.Display.DrawTetris(g.cur, false)
}

// NewTetris generate tetris at the top
func (g *Game) NewTetris() *Tetris {
	seed := rand.Intn(len(shapes))
	return g.newTetris(seed)
}

func (g *Game) newTetris(seed int) *Tetris {
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

func (g *Game) lineIsFull(line int) bool {
	for i := 0; i < g.Width; i++ {
		if g.board[line*g.Width+i] == nil {
			return false
		}
	}
	return true
}

func (g *Game) clearLines(base, lines int) int {
	if lines == 0 {
		return 0
	}

	for r := base; r >= lines; r-- {
		for c := 0; c < g.Width; c++ {
			g.board[r*g.Width+c] = g.board[(r-lines)*g.Width+c]
		}
	}

	for i := 0; i < g.Width*lines; i++ {
		g.board[i] = nil
	}

	// the more lines, the more score
	return []int{0, 1, 3, 5, 8}[lines]
}

func (g *Game) cleanUp() int {
	score := 0
	for row := g.Height - 1; row >= 0; {
		// skip the non-full row
		for row >= 0 && !g.lineIsFull(row) {
			row--
		}

		base, lines := row, 0
		// count the full row
		for row >= 0 && g.lineIsFull(row) {
			lines++
			row--
		}

		if lines > 0 {
			score += g.clearLines(base, lines)
		}
	}

	g.Display.DrawBlocks(g.board)
	return score
}

func (g *Game) persistTetris(t *Tetris) {
	// fill up the board
	s, o := t.shape, t.pos
	for i := 0; i < s.n*s.n; i++ {
		if s.data[i] != '.' {
			r, c := i/s.n, i%s.n
			g.board[(o.y+r)*g.Width+o.x+c] = s.block
		}
	}

	// clean up from bottom to up
	score := g.cleanUp()
	g.AddScore(score)
}

func (g *Game) AddScore(score int) {
	g.score <- score
}

func (g *Game) Down() {
	next := g.cur.Move(1, 0)

	// the tetris can not be move down any more,
	// then we persist the shape into the blocks
	// , clear the tetris and generate new one
	if g.isCollision(next) {
		// clear the old tetris
		g.Display.DrawTetris(g.cur, true)
		// persist to board
		g.persistTetris(g.cur)
		// re-draw blocks
		g.Display.DrawBlocks(g.board)

		// generate new tetris
		// if the new tetris is out of board then game is over
		g.cur, g.prev = g.NewTetris(), nil
		g.FlushTetris()

		if g.isCollision(g.cur) {
			close(g.over)
		}
	} else {
		g.cur, g.prev = next, g.cur
		g.FlushTetris()
	}
}

// Listen to the input and ticks
func (g *Game) Listen() {
	input := g.Input.Input()
	over := false
	for !over {
		select {
		case key := <-input:
			switch key {
			case Left:
				g.Transform(g.cur.Move(0, -1))
			case Right:
				g.Transform(g.cur.Move(0, 1))
			case Down:
				g.Down()
			case TurnRight:
				g.Transform(g.cur.Rotate())
			}
		case <-g.ticker.C:
			g.Down()
		case score := <-g.score:
			oldLevel := g.Level
			g.Score += score
			g.Level = g.Score / 1
			g.Display.UpdateScore(g.Score, g.Level)
			if g.Level > oldLevel {
				g.ticker.Stop()
				g.ticker = NewTicker(time.Duration(float64(g.ticker.Duration) * 0.95))
			}
		case <-g.over:
			over = true
		}
	}
}

func (g *Game) Init() {
	g.Input.Init()
	g.Display.Init()

	// init board
	g.board = make([]*Block, g.Width*g.Height)

	// init tetris
	g.cur = g.NewTetris()
	g.FlushTetris()

	g.AddScore(0)
}

func main() {
	width, height := 10, 10
	g := Game{
		Width:  width,
		Height: height,
		Input:  &TTYInput{},
		Display: &TTYDisplay{
			Origin: Point{10, 10},
			Width:  width,
			Height: height,
		},
		over:   make(chan struct{}),
		ticker: NewTicker(1000 * time.Millisecond),
		score:  make(chan int, 2),
	}

	g.Init()
	g.Listen()
	log.Println("game over")
}
