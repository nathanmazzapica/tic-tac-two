package dto

import (
	"time"
)

type Event interface{ isEvent() }

type Joined struct {
	Slot int
	Mark int
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
	State int
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

type ValidMove struct {
	R    int `json:"r"`
	C    int `json:"c"`
	Mark int `json:"mark"`
}

func (ValidMove) isEvent() {}

type InvalidMove struct {
	R   int   `json:"r"`
	C   int   `json:"c"`
	Err error `json:"err"`
}

func (InvalidMove) isEvent() {}

type Forfeited struct{}

func (Forfeited) isEvent() {}

type GameOver struct {
	Method int
	Winner int
	Line   [][2]int
}

func (GameOver) isEvent() {}
