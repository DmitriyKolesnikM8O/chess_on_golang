package board

const (
	boardSize = 8
	BlackKing = "\u2654"
	WhiteKing = "\u265A"

	BlackQueen = "\u2655"
	WhiteQueen = "\u265B"

	BlackRook = "\u2656"
	WhiteRook = "\u265C"

	BlackBishop = "\u2657"
	WhiteBishop = "\u265D"

	BlackKnight = "\u2658"
	WhiteKnight = "\u265E"

	BlackPawn = "\u2659"
	WhitePawn = "\u265F"

	illegalMoveMessage        = "Illegal move! Please enter again."
	causingSelfInCheckMessage = "This move will cause yourself in check! Please enter again."
	MovesLimitCount           = 400
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
