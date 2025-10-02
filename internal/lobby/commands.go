package lobby

import "github.com/nathanmazzapica/tic-tac-two/internal/game"

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

type AddSub struct {
	ID string
	Ch chan Event
}

func (AddSub) isCommand() {}

type RemSub struct {
	ID string
}

func (RemSub) isCommand() {}
