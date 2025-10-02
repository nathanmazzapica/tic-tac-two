package ws

type CommandSink interface {
	Post(cmd any)
}
