package lobby

import "github.com/nathanmazzapica/tic-tac-two/internal/game"

type Lobby struct {
	PlayersBySlot [2]*game.Player
	PlayersByID   map[string]*game.Player
	Board         game.Board
	Id            string
	Turn          game.Mark
	State         game.GameStatus
	Outcome       game.Outcome
	//connsById		map[string]*ws.Conn
}
