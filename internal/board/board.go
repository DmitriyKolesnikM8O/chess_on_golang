package board

import (
	"bytes"
	"strings"

	"github.com/DmitriyKolesnikM8O/chess_on_golang/pkg/utils"
)

func NewBoard() *Board {
	board := new(Board)
	squares := make([][]Square, boardSize)

	for i := 0; i < boardSize; i++ {
		squares[i] = make([]Square, boardSize)
		for j := 0; j < boardSize; j++ {
			squares[i][j] = Square{
				row:   i,
				col:   j,
				piece: nil,
			}
		}
	}
	board.squares = squares
	return board

}

func (board *Board) Setup(testCase utils.TestCase) {
	for _, id := range testCase.InitialPositions {
		board.initPiece(id.Position, id.Sign)
	}

	board.blackCaptures = testCase.BlackCapture
	board.whiteCaptures = testCase.WhiteCapture
}

func (board *Board) initPiece(position string, sign string) {
	square := board.GetSquare(position)
	if square.hasPiece() {
		panic("initPiece() failed on the position: " + position)
	}

	piece := CreatePiece(sign, square.row, square.col)
	square.SetPiece(&piece)
}

func (board Board) GetSquare(position string) *Square {
	col := int(position[0] - 'a')
	row := ((int)(position[1]-'0'))*-1 + boardSize

	return &board.squares[row][col]
}

func (square Square) hasPiece() bool {
	return square.piece != nil
}

func CreatePiece(sign string, row int, col int) Piece {
	var team Team

	if sign == strings.ToUpper(sign) {
		team = Black
	} else {
		team = White
	}

	return Piece{
		sign: sign,
		row:  row,
		col:  col,
		team: team,
	}
}

func (square *Square) SetPiece(piece *Piece) {
	square.piece = piece
}

func (board Board) String() string {
	var buffer bytes.Buffer

	boardString := make([][]string, boardSize)
	for i := 0; i < boardSize; i++ {
		boardString[i] = make([]string, boardSize)
		for j := 0; j < boardSize; j++ {
			boardString[i][j] = board.squares[i][j].String()
		}
	}

	buffer.WriteString(utils.StringifyBoard(boardString))

	buffer.WriteString("White captures: [")
	for _, capturedSign := range board.whiteCaptures {
		if capturedSign != "" {
			buffer.WriteString(getPieceSymbol(capturedSign) + "")
		}
	}

	buffer.WriteString("] \n")

	buffer.WriteString("Black captures: [")
	for _, capturedSign := range board.blackCaptures {
		if capturedSign != "" {
			buffer.WriteString(getPieceSymbol(capturedSign) + "")
		}
	}

	buffer.WriteString("] \n")

	return buffer.String()

}

func (square Square) String() string {
	if square.piece == nil {
		return ""
	}

	return square.piece.String()
}
func (piece Piece) String() string {
	return getPieceSymbol(piece.sign)
}

func getPieceSymbol(sign string) string {
	switch sign {
	case "k":
		return WhiteKing
	case "K":
		return BlackKing
	case "q":
		return WhiteQueen
	case "Q":
		return BlackQueen
	case "r":
		return WhiteRook
	case "R":
		return BlackRook
	case "b":
		return WhiteBishop
	case "B":
		return BlackBishop
	case "n":
		return WhiteKnight
	case "N":
		return BlackKnight
	case "p":
		return WhitePawn
	case "P":
		return BlackPawn

	default:
		panic("Cannot Print Piece. Unknown Sign:" + sign)
	}
}
