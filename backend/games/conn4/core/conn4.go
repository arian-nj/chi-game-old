package conn4_core

import "slices"

const (
	BOARD_WIDTH  = 7
	BOARD_HEIGHT = 6
	WIN_SIZE     = 4
)

type Conn4Board struct {
	Board []Cell
}

type Cell int

const (
	Empty Cell = iota
	One
	Two
)

func NewConn4Board() *Conn4Board {
	connBoard := &Conn4Board{
		Board: make([]Cell, BOARD_WIDTH*BOARD_HEIGHT),
	}

	return connBoard
}
func (t *Conn4Board) rcToIndex(r, c int) int {
	return r*BOARD_WIDTH + c
}

func (t *Conn4Board) indexToRC(index int) (int, int) {
	r := index / BOARD_WIDTH
	c := index % BOARD_WIDTH
	return r, c
}

func (t *Conn4Board) GetCell(cellindex int) Cell {
	return t.Board[cellindex]
}

func (board *Conn4Board) IsMoveValid(column int) (bool, string) {
	if (column < 0) && (column >= BOARD_WIDTH) {
		return false, "خارج از محدوده چیکار داری میکنی؟"
	}
	if board.Board[board.rcToIndex(0, column)] != Empty {
		return false, "این ستون پره"
	}

	return true, ""
}
func (board *Conn4Board) DropPiece(column int, cellType Cell) (int, bool) {
	for row := BOARD_HEIGHT - 1; row >= 0; row-- {
		idx := board.rcToIndex(row, column)
		if board.Board[idx] == Empty {
			board.Board[idx] = cellType
			return idx, true
		}
	}

	return -1, false
}

func (board *Conn4Board) HasWon(index int) (bool, []int) {
	moveType := board.Board[index]
	if moveType == Empty {
		return false, nil
	}

	r, c := board.indexToRC(index)
	directions := [][2]int{
		{0, 1},  // Horizontal check
		{1, 0},  // Vertical check
		{1, 1},  // Diagonal (\)
		{1, -1}, // Anti-Diagonal (/)
	}

	for _, dir := range directions {
		// Start the line with the current index
		line := []int{index}

		// Forward direction
		for i := 1; i < WIN_SIZE; i++ {
			nr, nc := r+(dir[0]*i), c+(dir[1]*i)
			if nr >= 0 && nr < BOARD_HEIGHT && nc >= 0 && nc < BOARD_WIDTH {
				ni := board.rcToIndex(nr, nc)
				if board.Board[ni] == moveType {
					line = append(line, ni)
				} else {
					break
				}
			}
		}

		// Backward direction
		for i := 1; i < WIN_SIZE; i++ {
			nr, nc := r-(dir[0]*i), c-(dir[1]*i)
			if nr >= 0 && nr < BOARD_HEIGHT && nc >= 0 && nc < BOARD_WIDTH {
				ni := board.rcToIndex(nr, nc)
				if board.Board[ni] == moveType {
					line = append(line, ni)
				} else {
					break
				}
			}
		}

		// If enough in line → return winner
		if len(line) >= WIN_SIZE {
			return true, line
		}
	}

	return false, nil
}
func (board *Conn4Board) IsAnyCellEmpty() bool {
	return slices.Contains(board.Board, Empty)
}
