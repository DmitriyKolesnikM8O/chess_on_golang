package game

import (
	"bufio"

	. "github.com/DmitriyKolesnikM8O/chess_on_golang/internal/board"
)

type ChessGame struct {
	board       *Board
	movesCount  int
	currentTeam Team
	reader      bufio.Reader
}
