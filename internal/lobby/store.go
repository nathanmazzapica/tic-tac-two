package lobby

import "sync"

type Store struct {
	mu   sync.RWMutex
	byID map[string]*Lobby
}

func (s *Store) Get(id string) *Lobby {
	return s.byID[id]
}
