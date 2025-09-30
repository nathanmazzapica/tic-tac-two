package lobby

import (
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"time"
)

type Event interface{ isEvent() }

type Joined struct {
	Slot int
	Mark game.Mark
}

func (Joined) isEvent() {}

type Reconnected struct {
	Slot int
}

func (Reconnected) isEvent() {}

type JoinRejected struct {
	Reason string
}

func (JoinRejected) isEvent() {}

type LobbyStateChanged struct {
	State State
}

func (LobbyStateChanged) isEvent() {}

type Left struct {
	Slot int
}

func (Left) isEvent() {}

type Paused struct {
	MissingSlot int
	Deadline    time.Time
}

func (Paused) isEvent() {}

type Resumed struct{}

func (Resumed) isEvent() {}

type StateChanged struct {
}

func (StateChanged) isEvent() {}

type InvalidMove struct {
	R   int `json:"r"`
	C   int `json:"c"`
	Err error
}

func (InvalidMove) isEvent() {}

type Forfeited struct{}

func (Forfeited) isEvent() {}

type GameOver struct{}

func (GameOver) isEvent() {}
