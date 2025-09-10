package game

type State int

const (
	InProgress State = iota
	Won
	Draw
)

type ApplyResult struct {
	GameStatus State
	Winner     Mark
	Line       Line
}

func (r ApplyResult) Terminal() bool {
	return r.GameStatus != InProgress
}
