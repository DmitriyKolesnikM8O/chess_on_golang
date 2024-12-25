package board

const (
	boardSize = 8
)

type Board struct {
	squares       [][]Square
	whiteCaptures []string
	blackCaptures []string
}

type Square struct {
	row   int
	col   int
	piece *Piece
}

type Piece struct {
	sign string
	team Team
	row  int
	col  int
}

type Team int

const (
	Undecided Team = iota
	White     Team = iota
	Black     Team = iota
)
