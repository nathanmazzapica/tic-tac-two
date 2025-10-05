package lobby

import (
	"context"
	"errors"
	"github.com/nathanmazzapica/tic-tac-two/internal/dto"
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"reflect"
	"testing"
	"time"
)

func newLobbyForTest(t *testing.T) (*Lobby, *TestProbe, context.CancelFunc) {
	t.Helper()
	probe := &TestProbe{C: make(chan dto.Event, 32)}
	l := newLobby(WithNotifier(probe))
	ctx, cancel := context.WithCancel(context.Background())
	go l.Run(ctx)
	return l, probe, cancel
}

func newInProgressLobbyForTest(t *testing.T) (*Lobby, *TestProbe, context.CancelFunc) {
	t.Helper()
	probe := &TestProbe{C: make(chan dto.Event, 32)}
	l := newLobby(WithNotifier(probe), StartInProgress())
	ctx, cancel := context.WithCancel(context.Background())
	go l.Run(ctx)
	return l, probe, cancel
}

func expectNext[T any](t *testing.T, ch <-chan dto.Event, match func(T) bool) {
	t.Helper()

	wantType := reflect.TypeOf((*T)(nil)).Elem()

	select {
	case evt := <-ch:
		gotType := reflect.TypeOf(evt)
		e, ok := evt.(T)
		if !ok || !match(e) {
			t.Fatalf("expecting %v but got %v: %+v; mismatch", wantType, gotType, evt)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timed out waiting for event")
	}
}

func TestPlayerJoinsIdleLobby(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(dto.Join{PlayerID: "p1"})

	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})

	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})
}

func TestJoinFlow(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	t.Run("first join", func(t *testing.T) {
		l.Post(dto.Join{PlayerID: "p1"})
		expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
			return e.Slot == 0 && e.Mark == game.X
		})

		expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
			return e.State == WaitingForSecond
		})
	})

	t.Run("second join", func(t *testing.T) {
		l.Post(dto.Join{PlayerID: "p2"})
		expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
			return e.Slot == 1 && e.Mark == game.O
		})

		expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
			return e.State == InProgress
		})
	})
}

func TestGameFullRejection(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(dto.Join{PlayerID: "p1"})
	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})
	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})

	l.Post(dto.Join{PlayerID: "p2"})
	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 1 && e.Mark == game.O
	})
	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == InProgress
	})

	t.Run("reject third player", func(t *testing.T) {
		l.Post(dto.Join{PlayerID: "p3"})
		expectNext[dto.JoinRejected](t, p.C, func(e dto.JoinRejected) bool {
			return e.Reason == "AlreadyStarted"
		})
	})
}

func TestInProgressLeave(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(dto.Join{PlayerID: "p1"})
	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})
	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})

	l.Post(dto.Join{PlayerID: "p2"})
	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 1 && e.Mark == game.O
	})
	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == InProgress
	})

	l.Post(dto.Leave{PlayerID: "p2"})
	expectNext[dto.Left](t, p.C, func(e dto.Left) bool {
		return e.Slot == 1
	})

	expectNext[dto.Paused](t, p.C, func(e dto.Paused) bool {
		return e.MissingSlot == 1
	})
}

func TestReconnect(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(dto.Join{PlayerID: "p1"})
	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})
	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})

	l.Post(dto.Join{PlayerID: "p2"})
	expectNext[dto.Joined](t, p.C, func(e dto.Joined) bool {
		return e.Slot == 1 && e.Mark == game.O
	})
	expectNext[dto.LobbyStateChanged](t, p.C, func(e dto.LobbyStateChanged) bool {
		return e.State == InProgress
	})

	t.Run("reconnect", func(t *testing.T) {
		l.Post(dto.Leave{PlayerID: "p2"})
		expectNext[dto.Left](t, p.C, func(e dto.Left) bool {
			return e.Slot == 1
		})

		expectNext[dto.Paused](t, p.C, func(e dto.Paused) bool {
			return e.MissingSlot == 1
		})

		l.Post(dto.Join{PlayerID: "p2"})
		expectNext[dto.Reconnected](t, p.C, func(e dto.Reconnected) bool {
			return e.Slot == 1
		})

		expectNext[dto.Resumed](t, p.C, func(e dto.Resumed) bool {
			return true
		})

	})
}

func TestMove(t *testing.T) {
	l, p, cancel := newInProgressLobbyForTest(t)
	defer cancel()

	t.Run("out of turn move", func(t *testing.T) {
		l.Post(dto.Move{
			R:    0,
			C:    0,
			Mark: game.O,
		})

		expectNext[dto.InvalidMove](t, p.C, func(e dto.InvalidMove) bool {
			return errors.Is(e.Err, game.ErrWrongTurn)
		})
	})

	t.Run("initial move", func(t *testing.T) {
		l.Post(dto.Move{
			R:    0,
			C:    0,
			Mark: game.X,
		})
		expectNext[dto.ValidMove](t, p.C, func(e dto.ValidMove) bool {
			return e.Mark == game.X && e.R == 0 && e.C == 0
		})
	})

	t.Run("move occupied spot", func(t *testing.T) {
		l.Post(dto.Move{
			R:    0,
			C:    0,
			Mark: game.O,
		})

		expectNext[dto.InvalidMove](t, p.C, func(e dto.InvalidMove) bool {
			return errors.Is(e.Err, game.ErrOccupied)
		})
	})

}
