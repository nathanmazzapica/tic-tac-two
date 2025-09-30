package lobby

import (
	"github.com/nathanmazzapica/tic-tac-two/internal/game"
	"time"
)

type Slot struct {
	PlayerID          string
	Mark              game.Mark
	Connected         bool
	reconnectDeadline time.Time
}

func (s *Slot) Occupied() bool {
	return s.PlayerID != ""
}

// On WS open set connected=true; lastSeen=now; connRef; clear deadline
// On msg received bump lastSeen
// On disconnect: connected=false; connRef=nil; deadline=now+grace
// On deadline expiry forfeit game
