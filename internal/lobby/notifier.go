package lobby

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
