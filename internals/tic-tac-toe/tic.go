package tictactoe

type Cell int

const (
	Empty Cell = iota
	X
	O
)

type TicBoard struct {
	Board [3][3]Cell
}

func NewTicBoard() *TicBoard {
	return &TicBoard{
		Board: [3][3]Cell{
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
		},
	}
}

func (borard *TicBoard) PlayMove(r, c int, moveType Cell) (bool, string) {
	if borard.Board[r][c] != Empty {
		return false, "خالی نیست"
	}
	borard.Board[r][c] = moveType
	return true, ""
}

func (board *TicBoard) HasWon() bool {
	// Check for a win in each row
	for i := range 3 {
		if board.Board[i][0] != Empty && board.Board[i][0] == board.Board[i][1] && board.Board[i][1] == board.Board[i][2] {
			return true
		}
	}

	// Check for a win in each column
	for i := range 3 {
		if board.Board[0][i] != Empty && board.Board[0][i] == board.Board[1][i] && board.Board[1][i] == board.Board[2][i] {
			return true
		}
	}

	// Check for a win on the main diagonal (top-left to bottom-right)
	if board.Board[0][0] != Empty && board.Board[0][0] == board.Board[1][1] && board.Board[1][1] == board.Board[2][2] {
		return true
	}

	// Check for a win on the anti-diagonal (top-right to bottom-left)
	if board.Board[0][2] != Empty && board.Board[0][2] == board.Board[1][1] && board.Board[1][1] == board.Board[2][0] {
		return true
	}

	// If no winning condition is met, return false
	return false
}

func (board *TicBoard) IsAnyCellEmpty() bool {
	for _, row := range board.Board {
		for _, cell := range row {
			if cell == Empty {
				return true
			}
		}
	}
	return false
}
