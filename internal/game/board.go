package game

type Mark int

const (
	Empty Mark = iota
	X
	O
)

type Board struct {
	cells [3][3]Mark
	moves uint8
}
