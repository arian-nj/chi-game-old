package xo

import (
	"math"
	"math/rand"
)

const (
	MaxDepth     = 6
	RandomChance = .25
)

func FindBestMove(board *XoBoard, player Cell) (int, int) {
	// startTime := time.Now()
	if rand.Float64() < RandomChance {
		return RandomMove(*board)

	}

	r, c := BestMove(board, player)
	// fmt.Println(time.Since(startTime).Microseconds())
	// fmt.Println("How many ", howMany)
	return r, c
	// return RandomMove(*board, player)
}

func BestMove(board *XoBoard, playerMove Cell) (int, int) {
	bestIndex := -1
	bestScore := math.MinInt

	for i := range board.MaxCellSize * board.MaxCellSize {
		if board.Board[i] != Empty {
			continue
		}
		board.Board[i] = playerMove
		score := MinMax(i, playerMove.Flip(), board, math.MinInt, math.MaxInt, MaxDepth, false)
		if score > bestScore {
			bestIndex = i
			bestScore = score
		}
		board.Board[i] = Empty
	}
	r, c := bestIndex/board.MaxCellSize, bestIndex%board.MaxCellSize
	return r, c
}

var howMany = 0

func MinMax(index int, Move Cell, board *XoBoard, alpha, beta int, depth int, isMaximizing bool) int {
	if depth == 0 {
		return 0
	}
	howMany++
	hasWon := board.HasWon(index)
	if hasWon {
		if isMaximizing {
			return -10 - depth
		}
		return 10 + depth
	}

	if !board.IsAnyCellEmpty() {
		return 0
	}

	if isMaximizing {
		bestScore := math.MinInt
		for i := range board.MaxCellSize * board.MaxCellSize {
			if board.Board[i] != Empty {
				continue
			}

			board.Board[i] = Move
			score := MinMax(i, Move.Flip(), board, alpha, beta, depth-1, false)
			board.Board[i] = Empty
			bestScore = max(bestScore, score)
			alpha = max(alpha, score)
			if beta <= alpha {
				break
			}

		}
		return bestScore
	}

	worstScore := math.MaxInt
	for i := range board.MaxCellSize * board.MaxCellSize {
		if board.Board[i] != Empty {
			continue
		}
		board.Board[i] = Move
		score := MinMax(i, Move.Flip(), board, alpha, beta, depth-1, true)
		board.Board[i] = Empty
		worstScore = min(worstScore, score)
		beta = min(beta, score)
		if beta <= alpha {
			break
		}

	}
	return worstScore

}

func RandomMove(board XoBoard) (int, int) {
	var emptyCells [][2]int
	for r := range board.MaxCellSize {
		for c := range board.MaxCellSize {
			if board.GetCell(board.toCellIndex(r, c)) == Empty {
				emptyCells = append(emptyCells, [2]int{r, c})
			}
		}
	}

	if len(emptyCells) == 0 {
		return -1, -1 // No available move
	}

	randomIndex := rand.Intn(len(emptyCells))
	move := emptyCells[randomIndex]
	return move[0], move[1]
}
