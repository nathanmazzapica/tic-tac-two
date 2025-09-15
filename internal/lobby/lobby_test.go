package lobby

import (
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCase struct {
	name      string
	setup     func(*Lobby)
	WantErr   error
	WantState State
}

func TestLobby_PlayerConnect(t *testing.T) {
	tests := []TestCase{
		{
			name:      "first player connects",
			WantErr:   nil,
			WantState: WaitingForSecond,
		},
		{
			name: "second player connects",
			setup: func(l *Lobby) {
				l.Connect(&game.Player{})
			},
			WantErr:   nil,
			WantState: InProgress, // Do we want in progress immediately? Should we look for "NextState" instead? idk
		},
		{
			name: "third player connects",
			setup: func(l *Lobby) {
				l.Connect(&game.Player{})
				l.Connect(&game.Player{})
			},
			WantErr:   ErrGameFull,
			WantState: InProgress,
		},
	}
	for _, tc := range tests {
		l := Lobby{}
		if tc.setup != nil {
			tc.setup(&l)
		}
		err := l.Connect(&game.Player{})
		assert.Equal(t, tc.WantErr, err)
		assert.Equal(t, tc.WantState, l.State)
	}
}

/*
func TestLobby_PlayerDisconnect(t *testing.T) {
	tests := []TestCase{
		{
			name: "waiting player disconnects",
			setup: func(l *Lobby) {
				l.Connect()
			},
			WantErr:   nil,
			WantState: Idle,
		},
	}
	for _, tc := range tests {
		l := Lobby{}
		if tc.setup != nil {
			tc.setup(&l)
		}
		err := l.Disconnect()
		assert.Equal(t, tc.WantErr, err)
		assert.Equal(t, tc.WantState, l.State)
	}
}
*/
