package lobby

import "sync"

type Store struct {
	mu   sync.RWMutex
	byID map[string]*Lobby
}

func New(opts ...Option) *Lobby {
	l := &Lobby{
		commands: make(chan Command, 64),
		state:    Idle,
		n:        newFanoutNotifier(),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (s *Store) Get(id string) *Lobby {
	return s.byID[id]
}
