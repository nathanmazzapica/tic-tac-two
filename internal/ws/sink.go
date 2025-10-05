package ws

import (
	"fmt"
	"github.com/nathanmazzapica/tic-tac-two/internal/dto"
)

type CommandSink interface {
	Post(cmd *dto.Envelope)
}

type Sink struct {
	commands chan<- dto.Command
}

// NewSink creates a Sink with a write-only channel of commands
func NewSink(commands chan<- dto.Command) *Sink {
	return &Sink{commands: commands}
}

func (s *Sink) Post(cmd *dto.Envelope) {
	switch cmd.Type {
	case "join":
		fmt.Println("received join command")
		joinCmd := dto.Join{PlayerID: cmd.Data["player_id"].(string)}
		fmt.Printf("Join: %+v\n", joinCmd)
		fmt.Printf("chan: %v\n", s.commands)
		s.commands <- joinCmd
	}
}
