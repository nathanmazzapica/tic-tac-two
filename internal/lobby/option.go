package lobby

type Option func(*Lobby)

func WithNotifier(n notifier) Option { return func(l *Lobby) { l.n = n } }

func WithInboxSize(n int) Option { return func(l *Lobby) { l.commands = make(chan Command, n) } }
