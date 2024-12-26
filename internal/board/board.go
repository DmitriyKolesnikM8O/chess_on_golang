package board

import (
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
