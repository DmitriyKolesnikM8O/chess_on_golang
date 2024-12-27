package board

const (
	boardSize = 8
	WhiteKing = "\u2654"
	BlackKing = "\u265A"

	WhiteQueen = "\u2655"
	BlackQueen = "\u265B"

	WhiteRook = "\u2656"
	BlackRook = "\u265C"

	WhiteBishop = "\u2657"
	BlackBishop = "\u265D"

	WhiteKnight = "\u2658"
	BlackKnight = "\u265E"

	WhitePawn = "\u2659"
	BlackPawn = "\u265F"
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
