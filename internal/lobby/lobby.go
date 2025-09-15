package lobby

import (
	"errors"
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
)

type State int

const (
	Idle State = iota
	WaitingForSecond
	InProgress
	Terminal
	Closed // Not sure about keeping this one, maybe for cleanup?
)

var (
	ErrGameFull     = errors.New("game is full")
	ErrGameOver     = errors.New("game is over")
	ErrInvalidToken = errors.New("invalid token")
)

type Lobby struct {
	PlayersBySlot [2]*game.Player
	PlayersByID   map[string]*game.Player
	Board         game.Board
	Id            string
	Turn          game.Mark
	State         State
	//Outcome       game.Outcome
	//connsById		map[string]*ws.Conn
}

func (l *Lobby) Init() {}

func (l *Lobby) Start() {}

func (l *Lobby) Tick() {}

func (l *Lobby) Connect(p *game.Player) error {
	if l.isFull() {
		return ErrGameFull
	}

	switch l.State {
	case Idle:
		l.PlayersBySlot[0] = p
		l.State = WaitingForSecond
	case WaitingForSecond:
		l.PlayersBySlot[1] = p
		l.State = InProgress
	case InProgress:
		if l.isFull() {
			return ErrGameFull
		}
		// otherwise a player disconnected
		// we need to find out which player
		// and if player that left = player trying to join
		// if so, reconnect; otherwise reject with Err
	case Terminal, Closed:
		return ErrGameOver
	}

	return nil
}

func (l *Lobby) Disconnect(p game.Player) error {
	return nil
}

func (l *Lobby) playerCount() int {
	count := 0
	for _, player := range l.PlayersBySlot {
		if player != nil {
			count++
		}
	}
	return count
}

func (l *Lobby) isFull() bool {
	return l.playerCount() == 2
}
