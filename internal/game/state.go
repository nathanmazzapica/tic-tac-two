package game

type State int

const (
	Idle State = iota
	Waiting
	InProgress
	Terminal
	Closed
)

type Method string

const (
	ThreeInARow Method = "win"
	Drawn       Method = "draw"
	Forfeit     Method = "forfeit"
	Timeout     Method = "timeout"
)

type Outcome struct {
	Winner Mark
	Method Method
	Line   string
}
