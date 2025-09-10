package game

import (
	"github.com/stretchr/testify/assert"
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
		})
	}
}
