package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: Rewrite everything here ☹️

func TestApplyGoodMove(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	assert.Equal(t, true, board.ApplyMove(2, 2, O))
}

func TestApplyBadMove(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	assert.Equal(t, false, board.ApplyMove(4, 2, X))
}

func TestApplyOver(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	board.ApplyMove(2, 2, O)
	t.Log(board.cells)
	assert.Equal(t, false, board.ApplyMove(2, 2, X))
}

func TestRowWin(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	board.ApplyMove(0, 0, O)
	board.ApplyMove(0, 1, O)
	board.ApplyMove(0, 2, O)

	type outcome struct {
		Win  bool
		Mark Mark
	}

	expected := outcome{true, O}
	win, mark := board.CheckWinner()
	assert.Equal(t, expected, outcome{win, mark})
}

func TestColWin(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	board.ApplyMove(0, 0, O)
	board.ApplyMove(1, 0, O)
	board.ApplyMove(2, 0, O)

	type outcome struct {
		Win  bool
		Mark Mark
	}

	expected := outcome{true, O}
	win, mark := board.CheckWinner()
	assert.Equal(t, expected, outcome{win, mark})
}

func TestDiagWin(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	board.ApplyMove(0, 0, O)
	board.ApplyMove(1, 1, O)
	board.ApplyMove(2, 2, O)

	type outcome struct {
		Win  bool
		Mark Mark
	}

	expected := outcome{true, O}
	win, mark := board.CheckWinner()
	assert.Equal(t, expected, outcome{win, mark})
}

func TestODiagWin(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	board.ApplyMove(0, 2, O)
	board.ApplyMove(1, 1, O)
	board.ApplyMove(2, 0, O)

	type outcome struct {
		Win  bool
		Mark Mark
	}

	expected := outcome{true, O}
	win, mark := board.CheckWinner()
	assert.Equal(t, expected, outcome{win, mark})
}

func TestNoWin(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	board.ApplyMove(0, 0, O)
	board.ApplyMove(1, 1, O)
	board.ApplyMove(2, 1, O)

	type outcome struct {
		Win  bool
		Mark Mark
		Line [][2]int
	}

	expected := outcome{false, Empty}
	win, mark, line := board.CheckWinner()
	assert.Equal(t, expected, outcome{win, mark})
}

func TestEmptyBoardWin(t *testing.T) {
	board := Board{[3][3]Mark{}, 0}
	type outcome struct {
		Win  bool
		Mark Mark
	}

	expected := outcome{false, Empty}
	win, mark := board.CheckWinner()
	assert.Equal(t, expected, outcome{win, mark})
}
