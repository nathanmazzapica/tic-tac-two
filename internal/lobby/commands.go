package lobby

import "github.com/nathanmazzapica/tic-tac-two/internal/game"

// Join, Leave, Move, Forfeit, Tick

type Command interface{ isCommand() }

type Join struct {
	PlayerID string `json:"player_id"`
}

func (Join) isCommand() {}

type Leave struct {
	PlayerID string `json:"player_id"`
}

func (Leave) isCommand() {}

type Move struct {
	R    int       `json:"r"`
	C    int       `json:"c"`
	Mark game.Mark `json:"mark"`
}

func (Move) isCommand() {}

type Forfeit struct {
	PlayerID string `json:"player_id"`
}

func (Forfeit) isCommand() {}

type Tick struct{}

func (Tick) isCommand() {}
