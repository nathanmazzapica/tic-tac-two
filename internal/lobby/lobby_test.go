package lobby

import (
	"context"
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
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

func expectNext[T any](t *testing.T, ch <-chan Event, match func(T) bool) {
	t.Helper()
	select {
	case evt := <-ch:
		e, ok := evt.(T)
		if !ok || !match(e) {
			t.Fatalf("expecting %T but got %v; mismatch", evt, evt)
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
			return e.State == Terminal
		})
	})
}
