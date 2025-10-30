package ai

import (
	"github.com/awesohame/gogo/internal/engine"
)

// Bot is interface for ai player
type Bot interface {
	// returns  chosen move
	SelectMove(board *engine.Board, color engine.Color) engine.Move
}

// returns the opp color
func opponentColor(c engine.Color) engine.Color {
	if c == engine.Black {
		return engine.White
	}
	return engine.Black
}
