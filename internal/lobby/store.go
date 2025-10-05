package lobby

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

var cleanupCycleDuration = time.Duration(time.Minute * 5)

type Store struct {
	mu   sync.RWMutex
	byID map[string]*Lobby
}

func NewStore() *Store {
	store := &Store{
		mu:   sync.RWMutex{},
		byID: make(map[string]*Lobby),
	}

	go func() {
		time.Sleep(cleanupCycleDuration)
		for {
			for id, l := range store.byID {
				if l.state == Closed {
					err := store.Cleanup(id)
					if err != nil {
						log.Printf("failed to cleanup lobby: %s", err)
					}
				}
			}
		}
	}()

	return store
}

// New is the exported wrapper for creating a new lobby that also attaches an ID and adds it to the LobbyStore map
func (s *Store) New(opts ...Option) *Lobby {
	l := newLobby(opts...)
	l.ID = uuid.New().String()
	s.mu.Lock()
	s.byID[l.ID] = l
	s.mu.Unlock()
	return l
}

func (s *Store) Cleanup(id string) error {
	if _, ok := s.byID[id]; ok {
		delete(s.byID, id)
		return nil
	}
	return errors.New("not found")
}

func (s *Store) Get(id string) *Lobby {
	return s.byID[id]
}
