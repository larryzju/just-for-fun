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
	Height = 12
)

type Point struct {
	x, y int
}

func cleanScreen() {
	fmt.Printf("%s[s", escape)
}

func moveCursor(line, col int) {
	fmt.Printf("%s[%d;%dH", escape, line, col)
}

func drawBoard(topLeft Point, width, height int) {
	moveCursor(topLeft.x, topLeft.y)
	fmt.Printf("%c", BorderTopLeft)
	for i := 0; i < width; i++ {
		fmt.Printf("%c", BorderHorization)
	}
	fmt.Printf("%c", BorderTopRight)

	for i := 1; i < height; i++ {
		moveCursor(topLeft.y+i, topLeft.x)
		fmt.Printf("%c", BorderVertical)
		moveCursor(topLeft.y+i, topLeft.x+width+1)
		fmt.Printf("%c", BorderVertical)
	}

	moveCursor(topLeft.y+height, topLeft.x)
	fmt.Printf("%c", BorderBottomLeft)
	for i := 0; i < width; i++ {
		fmt.Printf("%c", BorderHorization)
	}
	fmt.Printf("%c", BorderBottomRight)
}

func main() {
	cleanScreen()
	drawBoard(Point{10, 10}, Width, Height)
	moveCursor(12, 12)
	fmt.Printf("hahahaa")
	moveCursor(80, 0)
}
