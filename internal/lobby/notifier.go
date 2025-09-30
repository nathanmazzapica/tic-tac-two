package lobby

import "fmt"

type notifier interface {
	Broadcast(Event)
}

type TestProbe struct {
	C chan Event
}

func newTestProbe() *TestProbe {
	return &TestProbe{C: make(chan Event, 16)}
}

func (p *TestProbe) Broadcast(e Event) {
	p.C <- e
}

type fanoutNotifier struct {
	subscribers map[string]chan Event
	bufferSize  int
}

func newFanoutNotifier() *fanoutNotifier {
	return &fanoutNotifier{subscribers: make(map[string]chan Event), bufferSize: 128}
}

func (n *fanoutNotifier) Broadcast(e Event) {
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

func (n *fanoutNotifier) Subscribe(id string) <-chan Event {
	ch := make(chan Event, n.bufferSize)
	n.subscribers[id] = ch
	return ch
}
