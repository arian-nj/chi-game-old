package xo

import "slices"

type XoBoard struct {
	Board       []Cell
	MaxCellSize int
	WinSize     int
}

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

func (t *XoBoard) toCellIndex(r, c int) int {
	return r*t.MaxCellSize + c
}

func (t *XoBoard) toCellRC(index int) (r, c int) {
	r = index / t.MaxCellSize
	c = index % t.MaxCellSize
	return r, c
}

func (t *XoBoard) SetCell(cellIndex int, cell Cell) {
	t.Board[cellIndex] = cell
}

func (t *XoBoard) GetCell(cellindex int) Cell {
	return t.Board[cellindex]
}

func NewTicBoard(maxBoardSize int, winSize int) *XoBoard {
	if winSize > int(maxBoardSize) {
		panic("wtf win size is bigger than max board size")
	}

	ticBoard := &XoBoard{
		MaxCellSize: maxBoardSize,
		WinSize:     winSize,

		Board: make([]Cell, maxBoardSize*maxBoardSize),
	}

	return ticBoard
}
func (board *XoBoard) DeepCopy() *XoBoard {
	newBoard := make([]Cell, len(board.Board))
	copy(newBoard, board.Board)
	return &XoBoard{
		Board:       newBoard,
		MaxCellSize: board.MaxCellSize,
		WinSize:     board.WinSize,
	}
}

func (board *XoBoard) IsMoveValid(cellIndex int, moveType Cell) (bool, string) {
	r, c := board.toCellRC(cellIndex)
	if (r < 0 || c < 0) && (r >= board.MaxCellSize || c >= board.MaxCellSize) {
		return false, "خارج از محدوده چیکار داری میکنی؟"
	}

	if board.GetCell(cellIndex) != Empty {
		return false, "خالی نیست"
	}
	return true, ""
}

// HasWon checks if the last move at (r, c) resulted in a win.
func (board *XoBoard) HasWon(index int) bool {
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
			if nr >= 0 && nr < board.MaxCellSize && nc >= 0 && nc < board.MaxCellSize && board.GetCell(board.toCellIndex(nr, nc)) == moveType {
				count++
			} else {
				break
			}
		}

		// Check in the opposite direction (e.g., up-left)
		for i := 1; i < board.WinSize; i++ {
			nr, nc := r-(dir[0]*i), c-(dir[1]*i)
			if nr >= 0 && nr < board.MaxCellSize && nc >= 0 && nc < board.MaxCellSize && board.GetCell(board.toCellIndex(nr, nc)) == moveType {
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
func (board *XoBoard) IsAnyCellEmpty() bool {
	return slices.Contains(board.Board, Empty)
}
