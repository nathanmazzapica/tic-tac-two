package lobby

import "github.com/nathanmazzapica/tic-tac-two/internal/dto"

type Option func(*Lobby)

func WithNotifier(n notifier) Option { return func(l *Lobby) { l.n = n } }

func WithInboxSize(n int) Option { return func(l *Lobby) { l.commands = make(chan dto.Command, n) } }

func StartInProgress() Option {
	return func(l *Lobby) {
		l.slots[0].PlayerID = "p1"
		l.slots[0].Connected = true

		l.slots[1].PlayerID = "p2"
		l.slots[1].Connected = true

		l.state = InProgress
	}
}
