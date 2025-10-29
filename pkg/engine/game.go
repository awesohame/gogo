package engine

import (
	eng "github.com/awesohame/gogo/internal/engine"
	"github.com/awesohame/gogo/internal/game"
)

// public API for managing a Go game
type Game struct {
	session *game.Session
}

// creates new Go game custom board size
func NewGame(size int) *Game {
	return &Game{session: game.NewSession(size)}
}

// applies a move to the current game
func (g *Game) MakeMove(move eng.Move) error {
	return g.session.MakeMove(move)
}

// undoes the last move
func (g *Game) Undo() error {
	return g.session.Undo()
}

// redoes the next move
func (g *Game) Redo() error {
	return g.session.Redo()
}

// passes the current player's turn
func (g *Game) Pass() error {
	return g.session.Pass()
}

// ends the game with the current player resigning
func (g *Game) Resign() {
	g.session.Resign()
}

// returns the current board state
func (g *Game) CurrentBoard() *eng.Board {
	return g.session.CurrentBoard()
}

// returns whose turn it is
func (g *Game) CurrentTurn() eng.Color {
	return g.session.CurrentTurn()
}

// returns whether the game has ended
func (g *Game) IsGameOver() bool {
	return g.session.IsGameOver()
}

// returns the current score
func (g *Game) GetScore() eng.Score {
	return g.session.GetScore()
}

// returns the score with komi
func (g *Game) GetScoreWithKomi(komi float64) (float64, float64, eng.Color) {
	return g.session.GetScoreWithKomi(komi)
}

// returns the number of moves made
func (g *Game) MoveCount() int {
	return g.session.MoveCount()
}

// returns whether undo is possible
func (g *Game) CanUndo() bool {
	return g.session.CanUndo()
}

// returns whether redo is possible
func (g *Game) CanRedo() bool {
	return g.session.CanRedo()
}

// returns the board size
func (g *Game) Size() int {
	return g.session.Size()
}

// creates a Move at the given coords (1-based)
func (g *Game) NewMove(x, y int, color eng.Color) eng.Move {
	board := g.session.CurrentBoard()
	return eng.Move{
		Point: board.ToPoint(x, y),
		Color: color,
	}
}
