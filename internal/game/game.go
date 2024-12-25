package game

import (
	"bufio"
	"os"

	"github.com/DmitriyKolesnikM8O/chess_on_golang/internal/board"
)

func New() ChessGame {
	return ChessGame{
		board:       board.NewBoard(),
		movesCount:  0,
		currentTeam: board.Undecided,
		reader:      *bufio.NewReader(os.Stdin),
	}
}
