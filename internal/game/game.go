package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	chessongolang "github.com/DmitriyKolesnikM8O/chess_on_golang"
	"github.com/DmitriyKolesnikM8O/chess_on_golang/internal/board"
	"github.com/DmitriyKolesnikM8O/chess_on_golang/pkg/utils"
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
		panic("")
	}

	game.board.Setup(testCase)
	return testCase
}

func (game ChessGame) PrintGameStatus() {
	fmt.Println(game.board.String())
}

func (game *ChessGame) Start() {
	game.SetupBoard(chessongolang.PATH)

	game.PrintGameStatus()

	for {
		game.changeTurn(true)
		game.printAvailableMovesInCheck()
		input := game.promtInput(game.reader)
		end := game.execute(input)
		if end {
			return
		}

	}
}

func (game *ChessGame) execute(command string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			game.changeTurn(false)
		}
	}()

	checkmate := game.board.Execute(command, game.currentTeam)

	if checkmate {
		game.endGameWithWinner(getTeamName(game.currentTeam), "Checkmate", command)
		return true
	}

	if game.IsTie() {
		game.endGameByTie(command)
		return true
	}

	game.printAction(command)
	game.printGameStatus()
	return false
}

func (game ChessGame) endGameByTie(lastCommand string) {
	game.printAction(lastCommand)
	game.printGameStatus()
	fmt.Println("Tie game.  Too many moves.")
}

func (game ChessGame) IsTie() bool {
	return game.movesCount >= board.MovesLimitCount
}

func (game ChessGame) endGameWithWinner(winnerPlayer string, reason interface{}, lastCommand string) {
	game.printAction(lastCommand)
	game.printGameStatus()
	fmt.Println()
	fmt.Println(winnerPlayer, "player wins. ", reason)
}

func (game ChessGame) printGameStatus() {
	fmt.Println(game.board.String())
}

func (game ChessGame) printAction(action string) {
	fmt.Println(getTeamName(game.currentTeam), " player action: ", action)
}

func (game ChessGame) promtInput(reader bufio.Reader) string {
	fmt.Print(getTeamName(game.currentTeam), "> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")

	return input
}

func (game *ChessGame) changeTurn(next bool) {
	if next {
		game.movesCount++
	} else {
		game.movesCount--
	}

	switch game.currentTeam {
	case board.Undecided:
		game.currentTeam = board.White
	case board.Black:
		game.currentTeam = board.White
	case board.White:
		game.currentTeam = board.Black
	}
}

func (game ChessGame) printAvailableMovesInCheck() {
	current := game.currentTeam

	if !game.board.InCheck(current) {
		return
	}

	fmt.Println(getTeamName(game.currentTeam) + " is in check")
	fmt.Println("Available moves:")
	available := game.board.GetAvailableMovesInCheck(current)
	for _, move := range available {
		fmt.Println(move)
	}
	fmt.Println()
}

func getTeamName(current board.Team) string {
	switch current {
	case board.White:
		return "WHITE Player"
	case board.Black:
		return "BLACK Player"

	default:
		return "Unknown Player"
	}
}
