package ai

import (
	"github.com/awesohame/gogo/internal/engine"
)

// provides position evaluation
type Evaluator interface {
	Evaluate(board *engine.Board, color engine.Color) float64
}

// uses basic score diff
type SimpleEvaluator struct{}

func (e *SimpleEvaluator) Evaluate(board *engine.Board, color engine.Color) float64 {
	blackScore, whiteScore, _ := board.CalculateScoreWithKomi(6.5)

	if color == engine.Black {
		return blackScore - whiteScore
	}
	return whiteScore - blackScore
}

// considers territory
type InfluenceEvaluator struct{}

func (e *InfluenceEvaluator) Evaluate(board *engine.Board, color engine.Color) float64 {
	blackInf := GetInfluenceScore(board, engine.Black)
	whiteInf := GetInfluenceScore(board, engine.White)

	if color == engine.Black {
		return blackInf - whiteInf
	}
	return whiteInf - blackInf
}

// combines multiple evaluation methods
type HybridEvaluator struct {
	ScoreWeight     float64
	InfluenceWeight float64
}

func NewHybridEvaluator() *HybridEvaluator {
	return &HybridEvaluator{
		ScoreWeight:     0.7,
		InfluenceWeight: 0.3,
	}
}

func (e *HybridEvaluator) Evaluate(board *engine.Board, color engine.Color) float64 {
	// score based
	blackScore, whiteScore, _ := board.CalculateScoreWithKomi(6.5)
	scoreDiff := blackScore - whiteScore

	// influence based
	blackInf := GetInfluenceScore(board, engine.Black)
	whiteInf := GetInfluenceScore(board, engine.White)
	infDiff := blackInf - whiteInf

	// combine with weights
	combined := e.ScoreWeight*scoreDiff + e.InfluenceWeight*infDiff

	if color == engine.Black {
		return combined
	}
	return -combined
}
