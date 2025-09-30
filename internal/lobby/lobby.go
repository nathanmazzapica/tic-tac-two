package lobby

import (
	"context"
	"errors"
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"log"
	"time"
)

type State int

const (
	Idle State = iota
	WaitingForSecond
	InProgress
	Terminal
	Closed // Not sure about keeping this one, maybe for cleanup?
)

const GRACE_PERIOD = time.Second * 90

var (
	ErrGameFull         = errors.New("game is full")
	ErrAlreadyConnected = errors.New("already connected")
	ErrGameOver         = errors.New("game is over")
	// ErrInvalidToken indicates a player attempted to join a InProgress lobby that they were not a member of
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidPlayer = errors.New("invalid player")
	ErrStaleVersion  = errors.New("stale version")
)

// Lobby stores players and a board and governs actions taken on the board. It then produces events for WS to consume.
type Lobby struct {
	slots    [2]Slot
	board    game.Board
	ID       string
	Turn     game.Mark
	state    State
	subs     map[string]chan Event
	commands chan Command
	n        notifier
}

func New(opts ...Option) *Lobby {
	l := &Lobby{
		commands: make(chan Command, 64),
		subs:     make(map[string]chan Event),
		state:    Idle,
		n:        newTestProbe(),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Lobby) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case cmd, ok := <-l.commands:
			if !ok {
				return nil
			}
			switch c := cmd.(type) {
			case Join:
				l.handleJoin(c)
			case Leave:
				l.handleLeave(c)
			}
		default:
			log.Println("unknown command")
		}
	}
}

func (l *Lobby) Post(cmd Command) { l.commands <- cmd }

func (l *Lobby) handleJoin(cmd Join) {
	switch l.state {
	case Idle:
		if l.slots[0].Occupied() {
			l.n.Broadcast(JoinRejected{Reason: "CorruptState"})
			return
		}

		l.slots[0] = Slot{PlayerID: cmd.PlayerID, Mark: game.X, Connected: true}
		l.state = WaitingForSecond
		l.n.Broadcast(Joined{Slot: 0, Mark: game.X})
		l.n.Broadcast(LobbyStateChanged{State: l.state})

	case WaitingForSecond:
		if !l.slots[0].Occupied() || l.slots[1].Occupied() {
			l.n.Broadcast(JoinRejected{Reason: "CorruptState"})
			return
		}

		// reject duplicate join
		if idx, ok := l.findByPlayer(cmd.PlayerID); ok {
			s := l.slots[idx]
			s.Connected = true
			s.reconnectDeadline = time.Time{}
			l.slots[idx] = s
			l.n.Broadcast(Reconnected{Slot: idx})
			return
		}

		// find free slot deterministically
		if !l.slots[1].Occupied() {
			l.slots[1] = Slot{PlayerID: cmd.PlayerID, Mark: game.O, Connected: true}
			l.state = InProgress
			l.n.Broadcast(Joined{Slot: 1, Mark: game.O})
			l.n.Broadcast(LobbyStateChanged{State: l.state})
			return
		}
		l.n.Broadcast(JoinRejected{Reason: "LobbyFull"})

	case InProgress:
		if idx, ok := l.findByPlayer(cmd.PlayerID); ok {
			s := l.slots[idx]
			s.Connected = true
			s.reconnectDeadline = time.Time{}
			l.slots[idx] = s
			l.n.Broadcast(Reconnected{Slot: idx})
			return
		}
		l.n.Broadcast(JoinRejected{Reason: "AlreadyStarted"})
	}

}

func (l *Lobby) handleLeave(cmd Leave) {
	switch l.state {
	case InProgress:
		s, ok := l.findByPlayer(cmd.PlayerID)
		if ok {
			l.slots[s].Connected = false
			l.slots[s].reconnectDeadline = time.Now().Add(time.Second * 45)
		}
	}
}

func (l *Lobby) findByPlayer(id string) (int, bool) {
	if l.slots[0].PlayerID == id {
		return 0, true
	}

	if l.slots[1].PlayerID == id {
		return 1, true
	}

	return -1, false
}
