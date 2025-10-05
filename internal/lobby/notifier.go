package lobby

import (
	"fmt"
	"github.com/nathanmazzapica/tic-tac-two/internal/dto"
)

type notifier interface {
	Broadcast(dto.Event)
	Add(id string, ch chan dto.Event)
	Remove(id string)
}

type TestProbe struct {
	C chan dto.Event
}

func newTestProbe() *TestProbe {
	return &TestProbe{C: make(chan dto.Event, 16)}
}

func (p *TestProbe) Broadcast(e dto.Event) {
	p.C <- e
}

func (p *TestProbe) Add(id string, ch chan dto.Event) {}
func (p *TestProbe) Remove(id string)                 {}

type fanoutNotifier struct {
	subscribers map[string]chan dto.Event
	bufferSize  int
}

func newFanoutNotifier() *fanoutNotifier {
	return &fanoutNotifier{subscribers: make(map[string]chan dto.Event), bufferSize: 128}
}

func (n *fanoutNotifier) Broadcast(e dto.Event) {
	for id, ch := range n.subscribers {
		select {
		case ch <- e:
		default:
			fmt.Printf("Dropping unresponsive sub: %s\n", id)
			close(ch)
			delete(n.subscribers, id)
		}
	}
}

func (n *fanoutNotifier) Add(id string, ch chan dto.Event) { n.subscribers[id] = ch }

func (n *fanoutNotifier) Remove(id string) {
	if ch, ok := n.subscribers[id]; ok {
		close(ch)
		delete(n.subscribers, id)
	}
}
