package game

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DmitriyKolesnikM8O/chess_on_golang/internal/board"
	"github.com/DmitriyKolesnikM8O/chess_on_golang/pkg/utils"
)

const (
	InitialBoard = "./chess_on_golang/playBook/initial.txt"
)

func New() ChessGame {
	return ChessGame{
		board:       board.NewBoard(),
		movesCount:  0,
		currentTeam: board.Undecided,
		reader:      *bufio.NewReader(os.Stdin),
	}
}

func (game *ChessGame) SetupBoard(path string) utils.TestCase {
	testCase, err := utils.ParseTestCase(path)
	if err != nil {
		fmt.Printf("failed to parse test case: %v\n", err)
	}

	game.board.Setup(testCase)
	return testCase
}

func (game *ChessGame) Start() {
	game.SetupBoard(InitialBoard)
}
