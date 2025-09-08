package lobby

import "time"

type Session struct {
	playerID  string
	connected bool
	lastSeen  time.Time
	deadline  time.Time
	//connRef *ws.Conn
}

// On WS open set connected=true; lastSeen=now; connRef; clear deadline
// On msg received bump lastSeen
// On disconnect: connected=false; connRef=nil; deadline=now+grace
// On deadline expiry forfeit game
