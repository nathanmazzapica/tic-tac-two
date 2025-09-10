package game

type GameStatus int

const (
	InProgress GameStatus = iota
	Won
	Draw
)

type ApplyResult struct {
	GameStatus GameStatus
	Winner     Mark
	Line       Line
}

func (r ApplyResult) Terminal() bool {
	return r.GameStatus != InProgress
}
