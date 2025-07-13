package xo_core

type Cell int

const (
	Empty Cell = iota
	X
	O
)

func (c *Cell) Flip() Cell {
	switch *c {
	case X:
		return O
	case O:
		return X
	}
	return Empty

}

type TicBoard struct {
	Board       []Cell
	MaxCellSize int
	WinSize     int
}

func (t *TicBoard) CellIndex(r, c int) int {
	return r*t.MaxCellSize + c
}

func (t *TicBoard) SetCell(r, c int, cell Cell) {
	t.Board[t.CellIndex(r, c)] = cell
}

func (t *TicBoard) GetCell(r, c int) Cell {
	return t.Board[t.CellIndex(r, c)]
}

func NewTicBoard(maxBoardSize int, winSize int) *TicBoard {
	if winSize > int(maxBoardSize) {
		panic("wtf win size is bigger than max board size")
	}

	ticBoard := &TicBoard{
		MaxCellSize: maxBoardSize,
		WinSize:     winSize,

		Board: make([]Cell, maxBoardSize*maxBoardSize),
	}

	return ticBoard
}
func (board *TicBoard) DeepCopy() *TicBoard {
	newBoard := make([]Cell, len(board.Board))
	copy(newBoard, board.Board)
	return &TicBoard{
		Board:       newBoard,
		MaxCellSize: board.MaxCellSize,
		WinSize:     board.WinSize,
	}
}

func (board *TicBoard) PlayMove(r, c int, moveType Cell) (bool, string) {
	if (r < 0 || c < 0) && (r >= board.MaxCellSize || c >= board.MaxCellSize) {
		return false, "خارج از محدوده چیکار داری میکنی؟"
	}

	if board.GetCell(r, c) != Empty {
		return false, "خالی نیست"
	}
	board.SetCell(r, c, moveType)
	return true, ""
}

// HasWon checks if the last move at (r, c) resulted in a win.
func (board *TicBoard) HasWon(index int) bool {
	moveType := board.Board[index]
	r := index / board.MaxCellSize
	c := index % board.MaxCellSize

	// Each inner array is a {row_change, column_change} vector
	directions := [][2]int{
		{0, 1},  // Horizontal check
		{1, 0},  // Vertical check
		{1, 1},  // DIAGONAL check (\)
		{1, -1}, // ANTI-DIAGONAL check (/)
	}

	for _, dir := range directions {
		// Start at 1 to count the piece just played at (r, c)
		count := 1

		// Check in one direction (e.g., down-right)
		for i := 1; i < board.WinSize; i++ {
			nr, nc := r+(dir[0]*i), c+(dir[1]*i)
			if nr >= 0 && nr < board.MaxCellSize && nc >= 0 && nc < board.MaxCellSize && board.GetCell(nr, nc) == moveType {
				count++
			} else {
				break
			}
		}

		// Check in the opposite direction (e.g., up-left)
		for i := 1; i < board.WinSize; i++ {
			nr, nc := r-(dir[0]*i), c-(dir[1]*i)
			if nr >= 0 && nr < board.MaxCellSize && nc >= 0 && nc < board.MaxCellSize && board.GetCell(nr, nc) == moveType {
				count++
			} else {
				break
			}
		}

		// If the total count on this axis meets the win condition, we have a winner
		if count >= board.WinSize {
			return true
		}
	}

	// No winning line was found
	return false
}
func (board *TicBoard) IsAnyCellEmpty() bool {
	for _, cell := range board.Board {
		if cell == Empty {
			return true
		}
	}
	return false
}

// package tictactoe
//
// type Cell int
//
// const (
// 	Empty Cell = iota
// 	X
// 	O
// )
//
// type TicBoard struct {
// 	Board       [][]Cell
// 	MaxCellSize int
// 	WinSize     int
// }
//
// func NewTicBoard(maxBoardSize int, winSize int) *TicBoard {
// 	if winSize > int(maxBoardSize) {
// 		panic("wtf win size is bigger than max board size")
// 	}
// 	maxSize := int(maxBoardSize)
//
// 	board := make([][]Cell, maxSize)
// 	for r := range maxSize {
// 		board[r] = make([]Cell, maxSize)
// 		for c := range maxSize {
// 			board[r][c] = Empty
// 		}
// 	}
//
// 	return &TicBoard{
// 		Board:       board,
// 		MaxCellSize: maxSize,
// 		WinSize:     winSize,
// 	}
// }
//
// func (board *TicBoard) PlayMove(r, c int, moveType Cell) (bool, string) {
// 	if (r < 0 || c < 0) && (r >= board.MaxCellSize || c >= board.MaxCellSize) {
// 		return false, "خارج از محدوده چیکار داری میکنی؟"
// 	}
//
// 	if board.Board[r][c] != Empty {
// 		return false, "خالی نیست"
// 	}
// 	board.Board[r][c] = moveType
// 	return true, ""
// }
//
// // HasWon checks if the last move at (r, c) resulted in a win.
// func (board *TicBoard) HasWon(r, c int, moveType Cell) bool {
// 	// Each inner array is a {row_change, column_change} vector
// 	directions := [][2]int{
// 		{0, 1},  // Horizontal check
// 		{1, 0},  // Vertical check
// 		{1, 1},  // DIAGONAL check (\)
// 		{1, -1}, // ANTI-DIAGONAL check (/)
// 	}
//
// 	for _, dir := range directions {
// 		// Start at 1 to count the piece just played at (r, c)
// 		count := 1
//
// 		// Check in one direction (e.g., down-right)
// 		for i := 1; i < board.WinSize; i++ {
// 			nr, nc := r+(dir[0]*i), c+(dir[1]*i)
// 			if nr >= 0 && nr < board.MaxCellSize && nc >= 0 && nc < board.MaxCellSize && board.Board[nr][nc] == moveType {
// 				count++
// 			} else {
// 				break
// 			}
// 		}
//
// 		// Check in the opposite direction (e.g., up-left)
// 		for i := 1; i < board.WinSize; i++ {
// 			nr, nc := r-(dir[0]*i), c-(dir[1]*i)
// 			if nr >= 0 && nr < board.MaxCellSize && nc >= 0 && nc < board.MaxCellSize && board.Board[nr][nc] == moveType {
// 				count++
// 			} else {
// 				break
// 			}
// 		}
//
// 		// If the total count on this axis meets the win condition, we have a winner
// 		if count >= board.WinSize {
// 			return true
// 		}
// 	}
//
// 	// No winning line was found
// 	return false
// }
// func (board *TicBoard) IsAnyCellEmpty() bool {
// 	for _, row := range board.Board {
// 		for _, cell := range row {
// 			if cell == Empty {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
