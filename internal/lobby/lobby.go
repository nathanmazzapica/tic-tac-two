package lobby

import (
	"context"
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
)

const GracePeriod = time.Second * 90

// Lobby stores players and a board and governs actions taken on the board. It then produces events for WS to consume.
type Lobby struct {
	slots    [2]Slot
	board    game.Board
	ID       string
	Turn     game.Mark
	state    State
	commands chan Command
	n        notifier
}

// newLobby is the internal function for creating a new lobby with specified options
func newLobby(opts ...Option) *Lobby {
	l := &Lobby{
		commands: make(chan Command, 64),
		state:    Idle,
		n:        newFanoutNotifier(),
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
			case AddSub:
				l.n.Add(c.ID, c.Ch)
			case RemSub:
				l.n.Remove(c.ID)
			case Join:
				l.handleJoin(c)
			case Leave:
				l.handleLeave(c)
			case Move:
				l.handleMove(c)
			}
		default:
			log.Println("unknown command")
		}
	}
}

func (l *Lobby) Post(cmd Command) { l.commands <- cmd }

func (l *Lobby) Subscribe(id string) <-chan Event {
	ch := make(chan Event, 32)
	l.Post(AddSub{ID: id, Ch: ch})
	return ch
}

func (l *Lobby) Unsubscribe(id string) {
	l.Post(RemSub{ID: id})
}

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
			l.n.Broadcast(Resumed{})
			return
		}
		l.n.Broadcast(JoinRejected{Reason: "AlreadyStarted"})
	case Terminal:
		l.n.Broadcast(JoinRejected{Reason: "GameOver"})
	}

}

func (l *Lobby) handleLeave(cmd Leave) {
	switch l.state {
	case Idle:
		// this shouldn't happen, but i'll guard against... just in case...
		return
	case WaitingForSecond:
		l.state = Idle
		l.slots[0].PlayerID = ""
		// TODO: Set TTL
		l.n.Broadcast(LobbyStateChanged{l.state})
	case InProgress:
		s, ok := l.findByPlayer(cmd.PlayerID)
		if ok {
			reconnectDeadline := time.Now().Add(GracePeriod)
			l.slots[s].Connected = false
			l.slots[s].reconnectDeadline = reconnectDeadline
			l.n.Broadcast(Left{Slot: s})
			l.n.Broadcast(Paused{
				MissingSlot: s,
				Deadline:    reconnectDeadline,
			})
		}
	case Terminal:
		s, ok := l.findByPlayer(cmd.PlayerID)
		if ok {
			l.n.Broadcast(Left{Slot: s})
		}
	}
}

func (l *Lobby) handleMove(c Move) {
	switch l.state {
	case InProgress:
		row := c.R
		col := c.C
		mark := c.Mark
		res, err := l.board.ApplyMove(row, col, mark)
		if err != nil {
			// move was invalid
			l.n.Broadcast(InvalidMove{
				R:   row,
				C:   col,
				Err: err,
			})
			return
		}

		// We always need to broadcast move
		l.n.Broadcast(ValidMove{
			R:    row,
			C:    col,
			Mark: mark,
		})

		switch res.GameStatus {
		case game.Won, game.Draw:
			l.n.Broadcast(GameOver{
				Method: res.GameStatus,
				Winner: res.Winner,
				Line:   res.Line,
			})
		}
	default:
		// if we are anywhere but InProgress, reject move
		return
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
