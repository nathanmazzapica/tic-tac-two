package game

import "log"

type Mark int

const (
	Empty Mark = iota
	X
	O
)

var lines = [][][2]int{
	// rows
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	// col
	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	//diag
	{{0, 0}, {1, 1}, {2, 2}},
	{{0, 2}, {1, 1}, {2, 0}},
}

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
	log.Println(b.cells)
	for _, line := range lines {
		win, mark := b.checkLine(line)
		if win {
			return win, mark
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

func (b *Board) checkLine(line [][2]int) (bool, Mark) {
	origin := line[0]
	mark := b.cells[origin[0]][origin[1]]

	log.Println(line)
	log.Println(origin)
	log.Println(mark)
	log.Println("----")

	if mark == Empty {
		return false, Empty
	}

	for _, point := range line {
		r := point[0]
		c := point[1]

		if b.cells[r][c] != mark {
			return false, Empty
		}
	}

	return true, mark
}
