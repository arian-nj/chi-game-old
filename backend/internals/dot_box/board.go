package dotbox

const MaxCellSize = 5

type Cell int

const (
	Empty Cell = iota
	Blue
	Red
)

type DotBoxBoard struct {
	Board [5][5]Cell
}

func NewDotBoard() *DotBoxBoard {
	return &DotBoxBoard{
		Board: [5][5]Cell{
			{Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty},
		},
	}
}

func (board *DotBoxBoard) HasEmptyCell() bool {
	for _, row := range board.Board {
		for _, c := range row {
			if c == Empty {
				return true
			}
		}
	}
	return false
}

func (board *DotBoxBoard) PlayMove(r, c int, moveType Cell) (bool, string, bool) {
	if (r < 0 || c < 0) && (r >= MaxCellSize || c >= MaxCellSize) {
		return false, "خارج از محدوده چیکار داری میکنی؟", false
	}

	if board.Board[r][c] != Empty {
		return false, "اینجا خالی نیست", false
	}
	isScore := board.Score(r, c)
	board.Board[r][c] = moveType
	return true, "", isScore
}

func (board *DotBoxBoard) Score(r, c int) bool {
	var startRow, endRow int = 0, MaxCellSize - 1
	var startCol, endCol int = 0, MaxCellSize - 1

	if r-1 >= 0 {
		startRow = r - 1
	}
	if r+1 <= MaxCellSize-1 {
		endRow = r + 1
	}

	if c-1 >= 0 {
		startCol = c - 1
	}
	if c+1 <= MaxCellSize-1 {
		endCol = c + 1
	}
	return board.detectScore(startRow, endRow, startCol, endCol, r, c)
}

// func main() {
// 	board := NewDotBoard()
// 	orgBoard := board.Board
// 	cpBoard := board.Board
// 	cpBoard[0][1] = Blue
// 	fmt.Println(orgBoard)
// 	fmt.Println(cpBoard)
// }

func (board *DotBoxBoard) detectScore(startRow, endRow int, startCol, endCol int, targetR, targetC int) bool {
	orgBoard := board.Board
	cpBoard := board.Board
	cpBoard[targetR][targetC] = Blue

	for r := startRow; r <= endRow; r++ {
		for c := startCol; c <= endCol; c++ {
			if c > MaxCellSize-2 || r > MaxCellSize-2 {
				continue
			}
			orgResult := orgBoard[r][c] * orgBoard[r][c+1] * orgBoard[r+1][c] * orgBoard[r+1][c+1]
			cpResult := cpBoard[r][c] * cpBoard[r][c+1] * cpBoard[r+1][c] * cpBoard[r+1][c+1]
			if orgResult == 0 && cpResult != 0 {
				return true
			}
		}
	}
	return false
}
