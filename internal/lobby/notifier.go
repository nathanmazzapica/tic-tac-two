package lobby

import "fmt"

type notifier interface {
	Broadcast(Event)
	Add(id string, ch chan Event)
	Remove(id string)
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

func (p *TestProbe) Add(id string, ch chan Event) {}
func (p *TestProbe) Remove(id string)             {}

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

func (n *fanoutNotifier) Add(id string, ch chan Event) { n.subscribers[id] = ch }

func (n *fanoutNotifier) Remove(id string) {
	if ch, ok := n.subscribers[id]; ok {
		close(ch)
		delete(n.subscribers, id)
	}
}
