package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyDisplay struct{}

func (d *dummyDisplay) Init()                            {}
func (d *dummyDisplay) DrawTetris(t *Tetris, clear bool) {}
func (d *dummyDisplay) DrawBlocks(blocks []*Block)       {}
func (d *dummyDisplay) UpdateScore(score, level int)     {}

func TestRules(t *testing.T) {
	g := &Game{
		Width:   4,
		Height:  4,
		board:   make([]*Block, 4*4),
		Display: &dummyDisplay{},
	}

	fillBoard := func(s string) {
		for i, c := range s {
			if c != 'x' {
				g.board[i] = nil
			} else {
				g.board[i] = defaultBlock
			}
		}
	}

	dumpBoard := func() string {
		bytes := make([]byte, len(g.board))
		for i, b := range g.board {
			if b != nil {
				bytes[i] = 'x'
			} else {
				bytes[i] = '.'
			}
		}
		return string(bytes)
	}

	fillBoard("x...xxxx..x.xxxx")
	assert.True(t, g.lineIsFull(3))
	assert.False(t, g.lineIsFull(2))
	assert.True(t, g.lineIsFull(1))
	assert.False(t, g.lineIsFull(0))
	assert.Equal(t, g.cleanUp(), 2)
	assert.Equal(t, dumpBoard(), "........x.....x.")

	fillBoard("xxxxxxxxxxxxxxxx")
	assert.Equal(t, g.cleanUp(), 8)
	assert.Equal(t, dumpBoard(), "................")

	fillBoard("................")
	assert.Equal(t, g.cleanUp(), 0)
	assert.Equal(t, dumpBoard(), "................")

	fillBoard("..x.xxxx..x...x.")
	assert.Equal(t, g.cleanUp(), 1)
	assert.Equal(t, dumpBoard(), "......x...x...x.")
}
