package ai

import (
	"math"

	"github.com/awesohame/gogo/internal/engine"
)

// returns a heuristic score for a position
// +ve for Black, neg for White
func EvaluatePosition(board *engine.Board) float64 {
	score := 0.0

	// simple material count with territory estimate
	blackScore, whiteScore, _ := board.CalculateScoreWithKomi(0)
	score = blackScore - whiteScore

	return score
}

// checks if a move is imp (captures, prevent atari, inc liberties)
func IsCriticalMove(board *engine.Board, move engine.Move) bool {
	// Apply move and check immediate tactical consequences
	newBoard, err := board.ApplyMove(move)
	if err != nil {
		return false
	}

	// Check if move captures stones
	// Check if move saves a group in atari
	// Check if move threatens opponent groups
	// (simplified check - can be expanded)

	_ = newBoard
	return false
}

// calcs territorial influence for a color
// higher scores indicate stronger control
func GetInfluenceScore(board *engine.Board, color engine.Color) float64 {
	influence := 0.0
	size := board.Size()

	// count stones and weighted empty points nearby
	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			point := board.ToPoint(x, y)
			c := board.At(x, y)

			if c == color {
				influence += 1.0

				// Add influence for nearby empty points
				for _, n := range board.Neighbors(point) {
					nx, ny := board.ToXY(n)
					if nx >= 1 && nx <= size && ny >= 1 && ny <= size {
						if board.At(nx, ny) == engine.Empty {
							influence += 0.3
						}
					}
				}
			}
		}
	}

	return influence
}

// gives bonus for playing in corners/edges early
func GetCornerProximityBonus(board *engine.Board, point engine.Point) float64 {
	x, y := board.ToXY(point)
	size := board.Size()

	// dist from nearest corner
	minDist := math.Min(
		math.Min(float64(x-1), float64(size-x)),
		math.Min(float64(y-1), float64(size-y)),
	)

	// higher bonus for corner/edge moves
	if minDist <= 3 {
		return (4.0 - minDist) * 0.1
	}

	return 0.0
}

// checks if a move fills own eye (bad move)
func IsEyeFillingMove(board *engine.Board, move engine.Move) bool {
	point := move.Point
	color := move.Color
	size := board.Size()

	// must be empty
	x, y := board.ToXY(point)
	if board.At(x, y) != engine.Empty {
		return false
	}

	// check if surrounded by friendly stones
	neighbors := board.Neighbors(point)
	friendlyCount := 0

	for _, n := range neighbors {
		nx, ny := board.ToXY(n)
		if nx < 1 || nx > size || ny < 1 || ny > size {
			continue
		}

		nColor := board.At(nx, ny)
		switch nColor {
		case color:
			friendlyCount++
		case engine.Empty:
			return false // not surrounded
		}
	}

	// likely an eye if 3 or more friendly neighbors
	return friendlyCount >= 3
}
