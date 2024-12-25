package board

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
