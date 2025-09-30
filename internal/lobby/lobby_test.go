package lobby

import (
	"context"
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"reflect"
	"testing"
	"time"
)

func newLobbyForTest(t *testing.T) (*Lobby, *TestProbe, context.CancelFunc) {
	t.Helper()
	probe := &TestProbe{C: make(chan Event, 32)}
	l := New(WithNotifier(probe))
	ctx, cancel := context.WithCancel(context.Background())
	go l.Run(ctx)
	return l, probe, cancel
}

func newInProgressLobbyForTest(t *testing.T) (*Lobby, *TestProbe, context.CancelFunc) {
	t.Helper()
	probe := &TestProbe{C: make(chan Event, 32)}
	l := New(WithNotifier(probe), StartInProgress())
	ctx, cancel := context.WithCancel(context.Background())
	go l.Run(ctx)
	return l, probe, cancel
}

func expectNext[T any](t *testing.T, ch <-chan Event, match func(T) bool) {
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

	l.Post(Join{PlayerID: "p1"})

	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})

	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})
}

func TestJoinFlow(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	t.Run("first join", func(t *testing.T) {
		l.Post(Join{PlayerID: "p1"})
		expectNext[Joined](t, p.C, func(e Joined) bool {
			return e.Slot == 0 && e.Mark == game.X
		})

		expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
			return e.State == WaitingForSecond
		})
	})

	t.Run("second join", func(t *testing.T) {
		l.Post(Join{PlayerID: "p2"})
		expectNext[Joined](t, p.C, func(e Joined) bool {
			return e.Slot == 1 && e.Mark == game.O
		})

		expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
			return e.State == InProgress
		})
	})
}

func TestGameFullRejection(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(Join{PlayerID: "p1"})
	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})
	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})

	l.Post(Join{PlayerID: "p2"})
	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 1 && e.Mark == game.O
	})
	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == InProgress
	})

	t.Run("reject third player", func(t *testing.T) {
		l.Post(Join{PlayerID: "p3"})
		expectNext[JoinRejected](t, p.C, func(e JoinRejected) bool {
			return e.Reason == "AlreadyStarted"
		})
	})
}

func TestInProgressLeave(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(Join{PlayerID: "p1"})
	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})
	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})

	l.Post(Join{PlayerID: "p2"})
	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 1 && e.Mark == game.O
	})
	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == InProgress
	})

	l.Post(Leave{PlayerID: "p2"})
	expectNext[Left](t, p.C, func(e Left) bool {
		return e.Slot == 1
	})

	expectNext[Paused](t, p.C, func(e Paused) bool {
		return e.MissingSlot == 1
	})
}

func TestReconnect(t *testing.T) {
	l, p, cancel := newLobbyForTest(t)
	defer cancel()

	l.Post(Join{PlayerID: "p1"})
	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 0 && e.Mark == game.X
	})
	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == WaitingForSecond
	})

	l.Post(Join{PlayerID: "p2"})
	expectNext[Joined](t, p.C, func(e Joined) bool {
		return e.Slot == 1 && e.Mark == game.O
	})
	expectNext[LobbyStateChanged](t, p.C, func(e LobbyStateChanged) bool {
		return e.State == InProgress
	})

	t.Run("reconnect", func(t *testing.T) {
		l.Post(Leave{PlayerID: "p2"})
		expectNext[Left](t, p.C, func(e Left) bool {
			return e.Slot == 1
		})

		expectNext[Paused](t, p.C, func(e Paused) bool {
			return e.MissingSlot == 1
		})

		l.Post(Join{PlayerID: "p2"})
		expectNext[Reconnected](t, p.C, func(e Reconnected) bool {
			return e.Slot == 1
		})

		expectNext[Resumed](t, p.C, func(e Resumed) bool {
			return true
		})

	})
}

func TestMove(t *testing.T) {
	l, p, cancel := newInProgressLobbyForTest(t)
	defer cancel()

	t.Run("initial move", func(t *testing.T) {
		l.Post(Move{
			R:    0,
			C:    0,
			Mark: game.X,
		})
		expectNext[ValidMove](t, p.C, func(e ValidMove) bool {
			return e.Mark == game.X && e.R == 0 && e.C == 0
		})
	})
}
