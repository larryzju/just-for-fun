package main

import (
	"fmt"
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
	Height = 25
)

func cleanScreen() {
	fmt.Printf("%s[s", escape)
}

func moveCursor(line, col int) {
	fmt.Printf("%s[%d;%dH", escape, line, col)
}

func main() {
	cleanScreen()
	moveCursor(10, 10)
	fmt.Printf("HAHA")
}
