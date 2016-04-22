package board

import (
	"fmt"
	"io"
	"strconv"
)

type Coord [2]int

type Cell struct {
	Covered  bool
	Flagged  bool
	Contents string
}

func genBoard() [][]Cell {
	return [][]Cell{
		{Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}},
		{Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "*"}, Cell{true, false, "0"}, Cell{true, false, "0"}},
		{Cell{true, false, "*"}, Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}},
		{Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "*"}, Cell{true, false, "0"}},
		{Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "0"}, Cell{true, false, "*"}, Cell{true, false, "0"}},
	}
}

type Board struct {
	layout [][]Cell
	Over   bool
}

func NewBoard() *Board {
	layout := genBoard()

	b := &Board{layout, false}
	b.setup()
	return b
}

func (b *Board) getAdjacent(x, y int) [][2]int {
	coords := [][2]int{
		{x - 1, y - 1}, // bl
		{x + 0, y - 1}, // bm
		{x + 1, y - 1}, // br
		{x + 1, y + 0}, // mr
		{x + 1, y + 1}, // tr
		{x + 0, y + 1}, // tm
		{x - 1, y + 1}, // tl
		{x - 1, y + 0}, // ml
	}

	return coords
}

func (b *Board) String() string {
	var buf string
	for _, row := range b.layout {
		for _, cell := range row {
			if !b.Over {
				if cell.Flagged {
					buf += " ⚑"
					continue
				} else if cell.Covered {
					buf += " #"
					continue
				}
			}

			buf += " " + cell.Contents
		}
		buf += "\n"
	}

	return buf
}

// setup will re-render board with newly uncovered cells
func (b *Board) setup() {
	for x, row := range b.layout {
		for y, cell := range row {
			if cell.Contents == "*" {
				adj := b.getAdjacent(x, y)
				for _, coord := range adj {
					if (coord[0] == -1 || coord[0] >= len(row)) || (coord[1] == -1 || coord[1] >= len(b.layout)) {
						continue
					}
					adjCell := &b.layout[coord[0]][coord[1]]
					if val, err := strconv.Atoi(adjCell.Contents); err == nil {
						adjCell.Contents = strconv.Itoa(val + 1)
					}
				}
			}
		}
	}
}

func (b *Board) uncover(x, y, filterX, filterY int) {
	cell := &b.layout[y][x]
	if cell.Flagged {
		return
	}

	cell.Covered = false

	if cell.Contents == "*" {
		b.Over = true
	}

	if cell.Contents == "0" {
		adj := b.getAdjacent(x, y)
		for _, coord := range adj {
			if coord[0] == filterX && coord[1] == filterY {
				continue
			}
			if (coord[0] == -1 || coord[0] >= len(b.layout[0])) || (coord[1] == -1 || coord[1] >= len(b.layout)) {
				continue
			}
			b.uncover(coord[0], coord[1], x, y)
		}
	}
}

func (b *Board) Size() (int, int) {
	return len(b.layout[0]), len(b.layout)
}

// Flag will set cell to flagged
func (b *Board) Flag(x, y int) *Board {
	b.layout[y][x].Flagged = true
	return b
}

// Select will uncover cell and mutate board state
func (b *Board) Select(x, y int) *Board {
	b.uncover(x, y, -1, -1)
	return b
}

func (b *Board) Render(w io.Writer) {
	fmt.Fprintf(w, "%s", b)
}
