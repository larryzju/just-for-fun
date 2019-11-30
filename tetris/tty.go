// +build linux darwin
// +build cgo
package main

// #include <termios.h>
// #include <stdio.h>
import "C"
import (
	"os"
	"fmt"
	"log"
	"syscall"
)

// TCGetLocalMode get tty local mode setting
func TCGetLocalMode(fd uintptr) (uint64, error) {
	ctermios := C.struct_termios{}
	if e := C.tcgetattr(C.int(fd), &ctermios); e != 0 {
		return 0, fmt.Errorf("%s", C.perror(C.CString("tcgetattr failed")))
	}
	return uint64(ctermios.c_lflag), nil
}

// TCSetLocalMode set tty local mode setting
func TCSetLocalMode(fd uintptr, mode uint64) error {
	ctermios := C.struct_termios{}
	if e := C.tcgetattr(C.int(fd), &ctermios); e != 0 {
		return fmt.Errorf("%s", C.perror(C.CString("tcgetattr failed")))
	}

	ctermios.c_lflag = C.ulong(mode)
	if e := C.tcsetattr(C.int(fd), C.TCSANOW, &ctermios); e != 0 {
		return fmt.Errorf("%s", C.perror(C.CString("tcsetattr failed")))
	}

	return nil
}

// TTYInput wraps input of the tty
type TTYInput struct{}

// Init the console
func(tty *TTYInput) Init() {
	mode, err := TCGetLocalMode(os.Stdin.Fd())
	if err != nil {
		log.Fatal(err)
	}

	mode &= ^uint64(syscall.ECHO | syscall.ICANON)
	log.Printf("%08X\n", mode)
	if err := TCSetLocalMode(os.Stdin.Fd(), mode); err != nil {
		log.Fatal(err)
	}
}

// ReadInput get input from keyboard
func(tty *TTYInput) Input() <-chan Key {
	ch := make(chan Key, 3)
	go func() {
		input := make([]byte, 4)
		for {
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
// TTYDisplay show tetris on console
type TTYDisplay struct {
	Origin Point
	Width int
	Height int
}

// Init the display
func(tty *TTYDisplay) Init() {
	cleanScreen()
	hideCursor()
	drawBoard(tty.Origin, tty.Width, tty.Height)
}

func cleanScreen() {
	fmt.Printf("%s[s", escape)
}

func hideCursor() {
	// fmt.Printf("\e[?25l")
}

func drawBoard(topLeft Point, width, height int) {
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

func (d *TTYDisplay) DrawTetris(t *Tetris, clear bool) {
	o, s := t.pos, t.shape
	for i := 0; i < s.n*s.n; i++ {
		r, c := i/s.n, i%s.n
		if s.data[i] != '.' {
			moveCursor(d.Origin.y+o.y+r, d.Origin.x+o.x+c)
			char := s.data[i]
			if clear {
				char = ' '
			}
			fmt.Printf("%c", char)
		}
	}
}

func (d *TTYDisplay) DrawBlocks(blocks []*Block) {
	for i, b := range blocks {
		r, c := i/d.Width, i%d.Width
		moveCursor(d.Origin.y+r, d.Origin.x+c)
		if b != nil {
			fmt.Printf("@")
		} else {
			fmt.Printf(" ")
		}
	}
}

func (d *TTYDisplay) UpdateScore(score, level int) {
	moveCursor(d.Origin.y+1, d.Origin.x+d.Width+3)
	fmt.Printf("Score:%05d", score)

	moveCursor(d.Origin.y+2, d.Origin.x+d.Width+3)
	fmt.Printf("Level:   %2d", level)
}