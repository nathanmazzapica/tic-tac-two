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

func (b *Board) ApplyMove(r, c int, mark Mark) bool {
	if b.moves == 9 || (r > 2 || r < 0 || c > 2 || c < 0) {
		return false
	}

	if b.cells[r][c] != Empty {
		return false
	}

	b.cells[r][c] = mark
	b.moves++
	return true
}

func (b *Board) CheckWinner() (bool, Mark) {
	for r := 0; r <= 2; r++ {
		for c := 0; c <= 2; c++ {
			if b.cells[r][c] != Empty {
				continue
			}

		}
	}

	return false, Empty
}

func (b *Board) Reset() {
	for rowIndex, row := range b.cells {
		for colIndex, _ := range row {
			b.cells[rowIndex][colIndex] = Empty
		}
	}
	b.moves = 0
}
