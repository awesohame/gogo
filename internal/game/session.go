package game

import (
	"errors"

	"github.com/awesohame/gogo/internal/engine"
)

// Session manages a single game lifecycle with history for undo/redo
type Session struct {
	history      []*engine.Board // all board states in order
	currentIndex int             // points to curr position in history
	currentTurn  engine.Color    // whose turn it is
	blackPassed  bool            // true if black passed on last move
	whitePassed  bool            // true if white passed on last move
	gameOver     bool            // true if game has ended
	size         int             // board size
}

// creates a new game session with the specified board size
func NewSession(size int) *Session {
	initialBoard := engine.NewBoard(size)
	return &Session{
		history:      []*engine.Board{initialBoard},
		currentIndex: 0,
		currentTurn:  engine.Black, // black starts first
		blackPassed:  false,
		whitePassed:  false,
		gameOver:     false,
		size:         size,
	}
}

// applies a move to the current board state
// returns error if move is illegal or not the correct player's turn
func (s *Session) MakeMove(move engine.Move) error {
	if s.gameOver {
		return errors.New("game is over")
	}

	// check if correct player's turn
	if move.Color != s.currentTurn {
		return errors.New("not your turn")
	}

	// get curr board
	currentBoard := s.history[s.currentIndex]

	// apply move using engine
	newBoard, err := currentBoard.ApplyMove(move)
	if err != nil {
		return err
	}

	// truncate any future history if we're not at the end in case of undoes
	s.history = s.history[:s.currentIndex+1]

	// add new board to history
	s.history = append(s.history, newBoard)
	s.currentIndex++

	// reset pass flags since a move was made
	s.blackPassed = false
	s.whitePassed = false

	// switch turn
	s.currentTurn = s.opponentColor()

	return nil
}

// user passes turn
func (s *Session) Pass() error {
	if s.gameOver {
		return errors.New("game is over")
	}

	// mark that curr player passed
	if s.currentTurn == engine.Black {
		s.blackPassed = true
	} else {
		s.whitePassed = true
	}

	// if both players pass, game ends
	if s.blackPassed && s.whitePassed {
		s.gameOver = true
	}

	// switch turn
	s.currentTurn = s.opponentColor()

	return nil
}

// user resigns
func (s *Session) Resign() {
	s.gameOver = true
}

// moves game state back
func (s *Session) Undo() error {
	if s.currentIndex <= 0 {
		return errors.New("Cannot undo: At start of game")
	}

	s.currentIndex--

	// switch turn back
	s.currentTurn = s.opponentColor()

	// reset game over state if we undo from end
	s.gameOver = false
	s.blackPassed = false
	s.whitePassed = false

	return nil
}

// moves game state forward
func (s *Session) Redo() error {
	if s.currentIndex >= len(s.history)-1 {
		return errors.New("Cannot redo: At end of history")
	}

	s.currentIndex++

	// switch turn forward
	s.currentTurn = s.opponentColor()

	return nil
}

// returns the current board state
func (s *Session) CurrentBoard() *engine.Board {
	return s.history[s.currentIndex]
}

// returns whose turn it is
func (s *Session) CurrentTurn() engine.Color {
	return s.currentTurn
}

// returns whether the game has ended
func (s *Session) IsGameOver() bool {
	return s.gameOver
}

// calcs and returns the curr score
func (s *Session) GetScore() engine.Score {
	return s.CurrentBoard().CalculateChineseScore()
}

// calcs score with komi
func (s *Session) GetScoreWithKomi(komi float64) (black float64, white float64, winner engine.Color) {
	return s.CurrentBoard().CalculateScoreWithKomi(komi)
}

// returns the no. of moves made in the game
func (s *Session) MoveCount() int {
	return s.currentIndex
}

// returns whether undo is possible
func (s *Session) CanUndo() bool {
	return s.currentIndex > 0
}

// returns whether redo is possible
func (s *Session) CanRedo() bool {
	return s.currentIndex < len(s.history)-1
}

// returns the opposite of the curr turn
func (s *Session) opponentColor() engine.Color {
	if s.currentTurn == engine.Black {
		return engine.White
	}
	return engine.Black
}

// returns board size
func (s *Session) Size() int {
	return s.size
}
