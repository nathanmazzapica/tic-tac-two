package game

import (
	"errors"
	"log"
)

var (
	ErrOccupied    = errors.New("cell already occupied")
	ErrOutOfBounds = errors.New("cell out of bounds")
	ErrGameOver    = errors.New("game already finished")
	ErrEmptyMove   = errors.New("empty move")
)

type Mark int

const (
	Empty Mark = iota
	X
	O
)

const boardSize = 3

type Line [][2]int

var lines = []Line{
	// rows
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	// col
	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	//diag
	{{0, 0}, {1, 1}, {2, 2}},
	{{0, 2}, {1, 1}, {2, 0}},
}

type Board struct {
	cells [3][3]Mark
	moves uint8
}

func (b *Board) ApplyMove(r, c int, m Mark) (ApplyResult, error) {
	if m == Empty {
		return ApplyResult{}, ErrEmptyMove
	}

	if !inBounds(r, c) {
		return ApplyResult{}, ErrOutOfBounds
	}

	if b.cells[r][c] != Empty {
		return ApplyResult{}, ErrOccupied
	}

	b.cells[r][c] = m
	b.moves++

	won, winner, line := b.CheckWinner()

	if won {
		return ApplyResult{
			GameStatus: Won,
			Winner:     winner,
			Line:       line,
		}, nil
	}

	if b.IsFull() {
		return ApplyResult{
			GameStatus: Draw,
			Winner:     Empty,
			Line:       nil,
		}, nil
	}

	return ApplyResult{}, nil
}

func (b *Board) State() State {
	won, _, _ := b.CheckWinner()
	if won {
		return Won
	}

	if b.IsFull() {
		return Draw
	}
	return InProgress
}

func (b *Board) CheckWinner() (bool, Mark, Line) {
	log.Println(b.cells)
	for _, line := range lines {
		win, mark := b.checkLine(line)
		if win {
			return win, mark, line
		}
	}

	return false, Empty, nil
}

func (b *Board) Reset() {
	for rowIndex, row := range b.cells {
		for colIndex, _ := range row {
			b.cells[rowIndex][colIndex] = Empty
		}
	}
	b.moves = 0
}

func (b *Board) IsFull() bool {
	return b.moves == 9
}

func (b *Board) Cell(r, c int) Mark {
	return b.cells[r][c]
}

func (b *Board) checkLine(line Line) (bool, Mark) {
	origin := line[0]
	mark := b.cells[origin[0]][origin[1]]

	log.Println(line)
	log.Println(origin)
	log.Println(mark)
	log.Println("----")

	if mark == Empty {
		return false, Empty
	}

	for _, point := range line {
		r := point[0]
		c := point[1]

		if b.cells[r][c] != mark {
			return false, Empty
		}
	}

	return true, mark
}

func inBounds(r, c int) bool {
	return r >= 0 && r < boardSize && c >= 0 && c < boardSize
}
