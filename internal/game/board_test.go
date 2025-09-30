package game

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBoard_ApplyMove(t *testing.T) {
	assert.Equal(t, true, true)
	tests := []struct {
		name      string
		setup     func(*Board)
		row, col  int
		mark      Mark
		wantErr   error
		wantLine  Line
		wantState State
	}{
		{
			name: "valid move",
			row:  0, col: 0, mark: X,
			wantErr:  nil,
			wantLine: nil, wantState: InProgress,
		},
		{
			name: "invalid move",
			row:  3, col: 0, mark: X,
			wantErr:   ErrOutOfBounds,
			wantState: InProgress,
			wantLine:  nil,
		},
		{
			name: "empty move",
			row:  0, col: 0, mark: Empty,
			wantErr:   ErrEmptyMove,
			wantState: InProgress,
			wantLine:  nil,
		},
		{
			name: "already occupied",
			row:  2, col: 2, mark: X,
			setup: func(b *Board) {
				_, _ = b.ApplyMove(2, 2, O)
			},
			wantErr:   ErrOccupied,
			wantState: InProgress,
			wantLine:  nil,
		},
		{
			name: "win game",
			row:  0, col: 2, mark: O,
			setup: func(b *Board) {
				_, _ = b.ApplyMove(0, 0, O)
				_, _ = b.ApplyMove(0, 1, O)
			},
			wantErr:   nil,
			wantState: Won,
			wantLine:  Line{{0, 0}, {0, 1}, {0, 2}},
		},
		{
			name: "draw game",
			row:  2, col: 2, mark: X,
			setup: func(b *Board) {
				_, _ = b.ApplyMove(0, 0, O)
				_, _ = b.ApplyMove(0, 1, X)
				_, _ = b.ApplyMove(0, 2, O)
				_, _ = b.ApplyMove(1, 0, X)
				_, _ = b.ApplyMove(1, 1, X)
				_, _ = b.ApplyMove(1, 2, O)
				_, _ = b.ApplyMove(2, 0, X)
				_, _ = b.ApplyMove(2, 1, O)
			},
			wantErr:   nil,
			wantState: Draw,
			wantLine:  nil,
		},
		{
			name: "game already over (win)",
			row:  0, col: 2, mark: X,
			setup: func(b *Board) {
				_, _ = b.ApplyMove(0, 0, O)
				_, _ = b.ApplyMove(0, 1, X)
				_, _ = b.ApplyMove(0, 2, O)
				_, _ = b.ApplyMove(1, 0, X)
				_, _ = b.ApplyMove(1, 1, X)
				_, _ = b.ApplyMove(1, 2, O)
				_, _ = b.ApplyMove(2, 0, X)
				_, _ = b.ApplyMove(2, 1, O)
				_, _ = b.ApplyMove(2, 2, O)
			},
			wantErr:   ErrGameOver,
			wantState: Won,
		},
		{
			name: "game already over (draw)",
			row:  0, col: 2, mark: X,
			setup: func(b *Board) {
				_, _ = b.ApplyMove(0, 0, O)
				_, _ = b.ApplyMove(0, 1, X)
				_, _ = b.ApplyMove(0, 2, O)
				_, _ = b.ApplyMove(1, 0, X)
				_, _ = b.ApplyMove(1, 1, X)
				_, _ = b.ApplyMove(1, 2, O)
				_, _ = b.ApplyMove(2, 0, X)
				_, _ = b.ApplyMove(2, 1, O)
				_, _ = b.ApplyMove(2, 2, X)
			},
			wantErr:   ErrGameOver,
			wantState: Draw,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := Board{}
			if tc.setup != nil {
				tc.setup(&b)
			}
			res, err := b.ApplyMove(tc.row, tc.col, tc.mark)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantState, res.GameStatus)
			assert.Equal(t, tc.wantLine, res.Line)
			log.Println(b.Moves)
		})
	}
}
