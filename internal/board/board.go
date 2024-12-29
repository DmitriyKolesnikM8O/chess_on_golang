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
	if col < 0 || col >= boardSize {
		return nil
	}
	row := ((int)(position[1]-'0'))*-1 + boardSize
	if row < 0 || row >= boardSize {
		return nil
	}
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

func (board *Board) String() string {
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
			buffer.WriteString(getPieceSymbol(capturedSign) + " ")
		}
	}

	buffer.WriteString("] \n")

	buffer.WriteString("Black captures: [")
	for _, capturedSign := range board.blackCaptures {
		if capturedSign != "" {
			buffer.WriteString(getPieceSymbol(capturedSign) + " ")
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

func (board *Board) Execute(command string, team Team) bool {
	// Handle "in check situation" first
	if board.InCheck(team) {
		validMoves := board.GetAvailableMovesInCheck(team)
		if !containsMove(validMoves, command) {
			panic(illegalMoveMessage)
		}
	}

	tokens := strings.Split(command, " ")
	if len(tokens) != 2 {
		panic(illegalMoveMessage)
	}

	// Check the movePiece first, panic if it's illegal movePiece on board
	move := board.checkMove(tokens[0], tokens[1], team)

	// Move Piece
	board.movePiece(move.piece, move.squareFrom, move.squareTo)

	//Promote
	//TODO: Handle promotion for Pawn

	// return if the opponent team is in checkmate
	return board.inCheckmate(getOpponentTeam(team))
}

func (board Board) inCheckmate(curTeam Team) bool {
	return board.InCheck(curTeam) && len(board.GetAvailableMovesInCheck(curTeam)) == 0
}

type Move struct {
	piece                *Piece
	squareFrom, squareTo *Square
}

func (board Board) checkMove(origin, destination string, team Team) Move {

	//Check input positions
	squareFrom := board.getSquare(origin)
	if squareFrom == nil {
		panic(illegalMoveMessage)
	}
	squareTo := board.getSquare(destination)
	if squareTo == nil {
		panic(illegalMoveMessage)
	}
	piece := squareFrom.GetPiece()
	if piece == nil || piece.team != team {
		panic(illegalMoveMessage)

	}

	//Check if the piece's movement is valid
	moves := getMoves(board, *piece)
	if !containsMove(moves, destination) {
		panic(illegalMoveMessage)

	}

	//Check if causing self in check (considered as invalid movePiece in current rule)
	if board.moveWillCauseSelfCheck(origin, destination, team) {
		panic(causingSelfInCheckMessage)
	}

	return Move{piece, squareFrom, squareTo}

}

func (board Board) GetAvailableMovesInCheck(current Team) []string {
	var moves []string

	moves = append(moves, board.getAvailableMovesByMovingKing(current)...)
	moves = append(moves, board.getAvailableMovesByMovingOtherPieces(current)...)

	return moves
}

func (board Board) getAvailableMovesByMovingOtherPieces(current Team) []string {
	var moves []string
	threatenKingPieces := board.getThreateningKingPieces(current)
	if len(threatenKingPieces) != 1 {
		return moves
	}

	threatenPiece := threatenKingPieces[0]
	threatenPosition := getCoordinatePosition(threatenPiece.row, threatenPiece.col)

	for _, ownPiece := range board.getAllPiece(current) {
		if isKing(ownPiece) {
			continue //king's move was already considered in method getAvailableMovesByMovingKing()
		}
		positionFrom := getCoordinatePosition(ownPiece.row, ownPiece.col)

		for _, positionTo := range getMoves(board, ownPiece) {
			if !board.moveWillCauseSelfCheck(positionFrom, positionTo, current) {
				moves = append(moves, positionFrom+" "+threatenPosition)
			}
		}
	}

	return moves
}

func (board Board) moveWillCauseSelfCheck(positionFrom, positionTo string, team Team) bool {

	squareFrom := board.getSquare(positionFrom)
	squareTo := board.getSquare(positionTo)

	pieceFrom := squareFrom.piece
	pieceTo := squareTo.piece

	//Move piece and check
	board.movePiece(pieceFrom, squareFrom, squareTo)
	selfInCheck := board.InCheck(team)

	//Move pieces back
	board.movePiece(pieceFrom, squareTo, squareFrom)
	squareTo.setPiece(pieceTo)

	return selfInCheck
}

func (board *Board) movePiece(piece *Piece, squareFrom, squareTo *Square) {

	capturedPiece := squareTo.piece
	if capturedPiece != nil {
		board.captured(*capturedPiece)
	}

	//Update squares on board
	squareFrom.setPiece(nil)
	squareTo.setPiece(piece)

	//Update Piece row and col
	piece.row = squareTo.row
	piece.col = squareTo.col

}

func (square *Square) setPiece(piece *Piece) {
	square.piece = piece
}

func (board *Board) captured(capturedPiece Piece) {
	team := getOpponentTeam(capturedPiece.team)
	sign := capturedPiece.sign //board will print the symbols from piece.sign
	switch team {
	case White:
		// panic("TEST")
		board.whiteCaptures = append(board.whiteCaptures, sign)
	case Black:
		// panic("TEST")
		board.blackCaptures = append(board.blackCaptures, sign)
	}
}

func isKing(piece Piece) bool {
	symbol := getPieceSymbol(piece.sign)
	return symbol == WhiteKing || symbol == BlackKing
}

func (board Board) getThreateningKingPieces(current Team) []Piece {
	var threatenPieces []Piece
	kingPosition := board.getKingPosition(current)
	for _, opponentPiece := range board.getAllPiece(getOpponentTeam(current)) {
		opponentMoves := getMoves(board, opponentPiece)
		if containsMove(opponentMoves, kingPosition) {
			threatenPieces = append(threatenPieces, opponentPiece)
		}
	}
	return threatenPieces
}

func (board Board) getAvailableMovesByMovingKing(current Team) []string {
	var moves []string
	kingPosition := board.getKingPosition(current)
	kingPiece := board.GetSquare(kingPosition).piece

	opponentsMoves := board.getReachablePositions(getOpponentTeam(current))

	for _, kingMove := range getMoves(board, *kingPiece) {
		if !containsMove(opponentsMoves, kingMove) {
			moves = append(moves, kingPosition+" "+kingMove)
		}
	}

	return moves
}

func containsMove(moves []string, move string) bool {
	for _, element := range moves {
		if move == element {
			return true
		}
	}
	return false
}

func (board Board) InCheck(current Team) bool {
	kingPosition := board.getKingPosition(current)

	opponentPosition := board.getReachablePositions(getOpponentTeam(current))

	for _, position := range opponentPosition {
		if kingPosition == position {
			return true
		}
	}

	return false
}

func (board Board) getReachablePositions(current Team) []string {
	var moves []string
	pieces := board.getAllPiece(current)
	for _, piece := range pieces {
		moves = append(moves, getMoves(board, piece)...)
	}

	return moves
}

func getMoves(board Board, piece Piece) []string {

	row := piece.row
	col := piece.col
	team := piece.team

	switch piece.String() {
	case WhiteKing, BlackKing:
		return getKingMoves(row, col, board, team)
	case WhiteQueen, BlackQueen:
		return getQueenMoves(row, col, board, team)

	case WhiteRook, BlackRook:
		return getRookMoves(row, col, board, team)

	case WhiteBishop, BlackBishop:
		return getBishopMovesAt(row, col, board, team)

	case WhiteKnight, BlackKnight:
		return getKnightMoves(row, col, board, team)

	case WhitePawn, BlackPawn:
		return getPawnMoves(row, col, board, team)

	default:
		panic("piece.String() hasn't been defined: " + piece.String())
	}
}

func getPawnMoves(row, col int, board Board, team Team) []string {

	var moves []string

	var firstMove bool

	oneStepRow := row
	twoStepsRow := row

	switch team {
	case White:
		oneStepRow--
		twoStepsRow -= 2
		firstMove = row == boardSize-2
	case Black:
		oneStepRow++
		twoStepsRow += 2
		firstMove = row == 1
	}

	//Get one step forwards positions
	position := getCoordinatePosition(oneStepRow, col)
	if board.isEmptyAt(position) {
		moves = append(moves, position)
	}

	//Get two step forwards positions if it's first move
	position = getCoordinatePosition(twoStepsRow, col)
	if firstMove && board.isEmptyAt(position) {
		moves = append(moves, position)
	}

	//Get Two Killing positions if there is enemy nearby to kill
	position = getCoordinatePosition(oneStepRow, col-1)
	if board.canMoveTo(position, team) {
		moves = append(moves, position)
	}
	position = getCoordinatePosition(oneStepRow, col+1)
	if board.canMoveTo(position, team) {
		moves = append(moves, position)
	}

	return moves
}

func getKnightMoves(row, col int, board Board, team Team) []string {
	var moves []string

	for i := -2; i <= 2; i++ {
		if i == 0 {
			continue
		}
		var j int
		//up
		if i == -2 || i == 2 {
			j = 1
		} else {
			j = 2
		}
		leftward := getCoordinatePosition(row+i, col-j)
		if board.canMoveTo(leftward, team) {
			moves = append(moves, leftward)
		}
		rightward := getCoordinatePosition(row+i, col+j)
		if board.canMoveTo(rightward, team) {
			moves = append(moves, rightward)
		}
	}

	return moves
}

func getQueenMoves(row, col int, board Board, team Team) []string {
	var moves []string
	moves = append(moves, getRookMoves(row, col, board, team)...)
	moves = append(moves, getBishopMovesAt(row, col, board, team)...)

	return moves
}

func getBishopMovesAt(row, col int, board Board, team Team) []string {

	var moves []string

	//Move Top Left
	for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		position := getCoordinatePosition(i, j)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	//Move Top Right
	for i, j := row-1, col+1; i >= 0 && j < boardSize; i, j = i-1, j+1 {
		position := getCoordinatePosition(i, j)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	//Move Bottom Left
	for i, j := row+1, col-1; i < boardSize && j >= 0; i, j = i+1, j-1 {
		position := getCoordinatePosition(i, j)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	//Move Bottom Right
	for i, j := row+1, col+1; i < boardSize && j < boardSize; i, j = i+1, j+1 {
		position := getCoordinatePosition(i, j)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	return moves
}

func getRookMoves(row, col int, board Board, team Team) []string {
	var moves []string

	//left
	for j := col - 1; j >= 0; j-- {
		position := getCoordinatePosition(row, j)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	//Moving Right
	for j := col + 1; j < boardSize; j++ {
		position := getCoordinatePosition(row, j)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	//Moving Up
	for i := row - 1; i >= 0; i-- {
		position := getCoordinatePosition(i, col)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	//Moving Down
	for i := row + 1; i < boardSize; i++ {
		position := getCoordinatePosition(i, col)
		if board.isEmptyAt(position) {
			moves = append(moves, position)
		} else if board.canMoveTo(position, team) {
			moves = append(moves, position)
			break
		} else {
			break
		}
	}

	return moves
}

func (board Board) isEmptyAt(position string) bool {
	square := board.getSquare(position)
	return square != nil && !square.hasPiece()
}

func (board Board) getSquare(position string) *Square {
	col := int(position[0] - 'a')
	if col < 0 || col >= boardSize {
		return nil
	}
	row := ((int)(position[1]-'0'))*-1 + boardSize
	if row < 0 || row >= boardSize {
		return nil
	}
	return &board.squares[row][col]
}

func getKingMoves(row, col int, board Board, team Team) []string {

	var moves []string

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			position := getCoordinatePosition(row+i, col+j)
			if board.canMoveTo(position, team) {
				moves = append(moves, position)
			}
		}
	}

	return moves
}

func (board Board) canMoveTo(position string, team Team) bool {
	square := board.GetSquare(position)
	if square == nil {
		return false
	}
	piece := square.GetPiece()
	return piece == nil || piece.team != team
}

func (board Board) getAllPiece(current Team) []Piece {
	var pieces []Piece
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			piece := board.squares[i][j].GetPiece()
			if piece != nil && piece.team == current {
				pieces = append(pieces, *piece)
			}
		}
	}

	return pieces
}

func getOpponentTeam(current Team) Team {
	if current == White {
		return Black
	}

	return White
}

func (board Board) getKingPosition(current Team) string {
	var kingSymbol string

	if current == White {
		kingSymbol = WhiteKing
	} else {
		kingSymbol = BlackKing
	}

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			square := board.squares[i][j]
			piece := square.GetPiece()
			if piece != nil && piece.String() == kingSymbol {
				return getSquarePosition(square)
			}
		}
	}

	panic("Cannot find king in the board")
}

func (square Square) GetPiece() *Piece {
	return square.piece
}

func getSquarePosition(square Square) string {
	return getCoordinatePosition(square.row, square.col)
}

func getCoordinatePosition(row, col int) string {
	return (string)('a'+col) + (string)('0'-row+boardSize)
}
