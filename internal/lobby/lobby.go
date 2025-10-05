package lobby

import (
	"context"
	"fmt"
	"github.com/nathanmazzapica/tic-tac-two/internal/dto"
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"time"
)

type State int

const (
	Idle State = iota
	WaitingForSecond
	InProgress
	Terminal
	Closed
)

const GracePeriod = time.Second * 90

// Lobby stores players and a board and governs actions taken on the board. It then produces events for WS to consume.
type Lobby struct {
	slots    [2]Slot
	board    game.Board
	ID       string
	Turn     game.Mark
	state    State
	commands chan dto.Command
	n        notifier
}

// newLobby is the internal function for creating a new lobby with specified options
func newLobby(opts ...Option) *Lobby {
	l := &Lobby{
		commands: make(chan dto.Command, 64),
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
			fmt.Println("yep")
			if !ok {
				return nil
			}
			switch c := cmd.(type) {
			case dto.AddSub:
				l.n.Add(c.ID, c.Ch)
			case dto.RemSub:
				l.n.Remove(c.ID)
			case dto.Join:
				l.handleJoin(c)
			case dto.Leave:
				l.handleLeave(c)
			case dto.Move:
				l.handleMove(c)
			default:
				fmt.Println("unknown command")
			}
		}
	}
}

func (l *Lobby) Post(cmd dto.Command) { l.commands <- cmd }

func (l *Lobby) Sink() chan<- dto.Command {
	return l.commands
}

func (l *Lobby) Subscribe(id string) <-chan dto.Event {
	ch := make(chan dto.Event, 32)
	l.Post(dto.AddSub{ID: id, Ch: ch})
	return ch
}

func (l *Lobby) Unsubscribe(id string) {
	l.Post(dto.RemSub{ID: id})
}

func (l *Lobby) handleJoin(cmd dto.Join) {
	fmt.Println("join!")
	switch l.state {
	case Idle:
		if l.slots[0].Occupied() {
			l.n.Broadcast(dto.JoinRejected{Reason: "CorruptState"})
			return
		}

		l.slots[0] = Slot{PlayerID: cmd.PlayerID, Mark: game.X, Connected: true}
		l.state = WaitingForSecond
		l.n.Broadcast(dto.Joined{Slot: 0, Mark: int(game.X)})
		l.n.Broadcast(dto.LobbyStateChanged{State: int(l.state)})

	case WaitingForSecond:
		if !l.slots[0].Occupied() || l.slots[1].Occupied() {
			l.n.Broadcast(dto.JoinRejected{Reason: "CorruptState"})
			return
		}

		// reject duplicate join
		if idx, ok := l.findByPlayer(cmd.PlayerID); ok {
			s := l.slots[idx]
			s.Connected = true
			s.reconnectDeadline = time.Time{}
			l.slots[idx] = s
			l.n.Broadcast(dto.Reconnected{Slot: idx})
			return
		}

		// find free slot deterministically
		if !l.slots[1].Occupied() {
			l.slots[1] = Slot{PlayerID: cmd.PlayerID, Mark: game.O, Connected: true}
			l.state = InProgress
			l.n.Broadcast(dto.Joined{Slot: 1, Mark: int(game.O)})
			l.n.Broadcast(dto.LobbyStateChanged{State: int(l.state)})
			return
		}
		l.n.Broadcast(dto.JoinRejected{Reason: "LobbyFull"})

	case InProgress:
		if idx, ok := l.findByPlayer(cmd.PlayerID); ok {
			s := l.slots[idx]
			s.Connected = true
			s.reconnectDeadline = time.Time{}
			l.slots[idx] = s
			l.n.Broadcast(dto.Reconnected{Slot: idx})
			l.n.Broadcast(dto.Resumed{})
			return
		}
		l.n.Broadcast(dto.JoinRejected{Reason: "AlreadyStarted"})
	case Terminal:
		l.n.Broadcast(dto.JoinRejected{Reason: "GameOver"})
	}

}

func (l *Lobby) handleLeave(cmd dto.Leave) {
	switch l.state {
	case Idle:
		// this shouldn't happen, but i'll guard against... just in case...
		return
	case WaitingForSecond:
		l.state = Idle
		l.slots[0].PlayerID = ""
		// TODO: Set TTL
		l.n.Broadcast(dto.LobbyStateChanged{int(l.state)})
	case InProgress:
		s, ok := l.findByPlayer(cmd.PlayerID)
		if ok {
			reconnectDeadline := time.Now().Add(GracePeriod)
			l.slots[s].Connected = false
			l.slots[s].reconnectDeadline = reconnectDeadline
			l.n.Broadcast(dto.Left{Slot: s})
			l.n.Broadcast(dto.Paused{
				MissingSlot: s,
				Deadline:    reconnectDeadline,
			})
		}
	case Terminal:
		s, ok := l.findByPlayer(cmd.PlayerID)
		if ok {
			l.n.Broadcast(dto.Left{Slot: s})
		}
	}
}

func (l *Lobby) handleMove(c dto.Move) {
	switch l.state {
	case InProgress:
		row := c.R
		col := c.C
		mark := c.Mark
		res, err := l.board.ApplyMove(row, col, game.Mark(mark))
		if err != nil {
			// move was invalid
			l.n.Broadcast(dto.InvalidMove{
				R:   row,
				C:   col,
				Err: err,
			})
			return
		}

		// We always need to broadcast move
		l.n.Broadcast(dto.ValidMove{
			R:    row,
			C:    col,
			Mark: int(mark),
		})

		switch res.GameStatus {
		case game.Won, game.Draw:
			l.n.Broadcast(dto.GameOver{
				Method: int(res.GameStatus),
				Winner: int(res.Winner),
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
