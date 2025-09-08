package lobby

import "sync"

type LobbyStore struct {
	mu   sync.RWMutex
	byID map[string]*Lobby
}
