package board

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"
)

const (
	CellNone = iota
	CellBomb
)

const Density float64 = 0.7

var CellTypes = []string{
	"0",
	"💣",
}

type Coord [2]int

type Cell struct {
	Covered  bool
	Flagged  bool
	Contents string
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func genLayout(width, height int) [][]Cell {
	layout := make([][]Cell, height)
	for y := 0; y < width; y++ {
		layout[y] = make([]Cell, width)
		for x := 0; x < height; x++ {
			firstNum := rand.Float64()

			var typ int
			if firstNum >= Density {
				typ = rand.Intn(len(CellTypes))
			} else {
				typ = CellNone
			}

			layout[y][x] = Cell{true, false, CellTypes[typ]}
		}
	}
	return layout
}

type Board struct {
	layout    [][]Cell
	winCount  int
	currCount int
	Over      bool
	Won       bool
}

func NewBoard() *Board {
	var x, y = 19, 19
	layout := genLayout(x, y)

	b := &Board{layout, x * y, 0, false, false}
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
					buf += " ▣"
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
			if cell.Contents == CellTypes[CellBomb] {
				adj := b.getAdjacent(x, y)
				for _, coord := range adj {
					if (coord[0] == -1 || coord[0] >= len(row)) ||
						(coord[1] == -1 || coord[1] >= len(b.layout)) {
						continue
					}
					adjCell := &b.layout[coord[0]][coord[1]]
					if val, err := strconv.Atoi(adjCell.Contents); err == nil {
						val += 1

						// number := color.New().SprintfFunc()
						// switch {
						// case val > 0:
						// 	number = color.New(color.FgGreen).SprintfFunc()
						// 	fallthrough
						// case val >= 2:
						// 	number = color.New(color.FgYellow).SprintfFunc()
						// 	fallthrough
						// case val >= 3:
						// 	number = color.New(color.FgRed).SprintfFunc()
						// }
						adjCell.Contents = fmt.Sprintf("%s", strconv.Itoa(val))
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

	if cell.Contents == CellTypes[CellBomb] {
		b.Over = true
		return
	}

	b.winCount++

	if cell.Contents == "0" {
		adj := b.getAdjacent(x, y)
		for _, coord := range adj {
			if (coord[0] == -1 || coord[0] >= len(b.layout[0])) ||
				(coord[1] == -1 || coord[1] >= len(b.layout)) {
				continue
			}
			if !b.layout[coord[1]][coord[0]].Covered {
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
	cell := &b.layout[y][x]
	if cell.Flagged {
		cell.Covered = true
		cell.Flagged = false
	} else {
		cell.Covered = false
		cell.Flagged = true
	}
	return b
}

// Select will uncover cell and mutate board state
func (b *Board) Select(x, y int) *Board {
	b.uncover(x, y, -1, -1)
	if b.currCount == b.winCount {
		b.Won = true
	}
	return b
}

func (b *Board) Render(w io.Writer) {
	fmt.Fprintf(w, "%s", b)
}
