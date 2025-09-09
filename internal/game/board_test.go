package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
