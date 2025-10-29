package tests

import (
	"testing"

	eng "github.com/awesohame/gogo/internal/engine"
	"github.com/awesohame/gogo/pkg/gogo/engine"
)

// TestGameCreation tests basic game initialization
func TestGameCreation(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"9x9 board", 9},
		{"13x13 board", 13},
		{"19x19 board", 19},
		{"5x5 board", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := engine.NewGame(tt.size)
			if game == nil {
				t.Fatal("Failed to create game")
			}

			if game.Size() != tt.size {
				t.Errorf("Expected size %d, got %d", tt.size, game.Size())
			}

			if game.CurrentTurn() != eng.Black {
				t.Errorf("Expected Black to start, got %v", game.CurrentTurn())
			}

			if game.IsGameOver() {
				t.Error("New game should not be over")
			}

			if game.MoveCount() != 0 {
				t.Errorf("New game should have 0 moves, got %d", game.MoveCount())
			}
		})
	}
}

// TestSimpleCapture tests basic capture mechanics
func TestSimpleCapture(t *testing.T) {
	game := engine.NewGame(9)

	// Black plays
	move1 := game.NewMove(4, 4, eng.Black)
	if err := game.MakeMove(move1); err != nil {
		t.Fatalf("Move 1 failed: %v", err)
	}

	// White plays next to black
	move2 := game.NewMove(5, 4, eng.White)
	if err := game.MakeMove(move2); err != nil {
		t.Fatalf("Move 2 failed: %v", err)
	}

	// Black surrounds white stone (part 1)
	move3 := game.NewMove(5, 3, eng.Black)
	if err := game.MakeMove(move3); err != nil {
		t.Fatalf("Move 3 failed: %v", err)
	}

	// White plays elsewhere
	move4 := game.NewMove(2, 2, eng.White)
	if err := game.MakeMove(move4); err != nil {
		t.Fatalf("Move 4 failed: %v", err)
	}

	// Black continues surrounding
	move5 := game.NewMove(5, 5, eng.Black)
	if err := game.MakeMove(move5); err != nil {
		t.Fatalf("Move 5 failed: %v", err)
	}

	// White plays elsewhere
	move6 := game.NewMove(3, 2, eng.White)
	if err := game.MakeMove(move6); err != nil {
		t.Fatalf("Move 6 failed: %v", err)
	}

	// Black captures white stone
	move7 := game.NewMove(6, 4, eng.Black)
	if err := game.MakeMove(move7); err != nil {
		t.Fatalf("Move 7 (capture) failed: %v", err)
	}

	// Verify white stone was captured
	board := game.CurrentBoard()
	if board.At(5, 4) != eng.Empty {
		t.Error("White stone should have been captured")
	}

	if game.MoveCount() != 7 {
		t.Errorf("Expected 7 moves, got %d", game.MoveCount())
	}
}

// TestSuicideRule tests that suicide moves are prevented
func TestSuicideRule(t *testing.T) {
	game := engine.NewGame(9)

	// Set up a position where suicide would occur
	// Black at (4,4)
	if err := game.MakeMove(game.NewMove(4, 4, eng.Black)); err != nil {
		t.Fatalf("Setup move 1 failed: %v", err)
	}

	// White at (3,4)
	if err := game.MakeMove(game.NewMove(3, 4, eng.White)); err != nil {
		t.Fatalf("Setup move 2 failed: %v", err)
	}

	// Black at (4,3)
	if err := game.MakeMove(game.NewMove(4, 3, eng.Black)); err != nil {
		t.Fatalf("Setup move 3 failed: %v", err)
	}

	// White at (5,4)
	if err := game.MakeMove(game.NewMove(5, 4, eng.White)); err != nil {
		t.Fatalf("Setup move 4 failed: %v", err)
	}

	// Black at (4,5)
	if err := game.MakeMove(game.NewMove(4, 5, eng.Black)); err != nil {
		t.Fatalf("Setup move 5 failed: %v", err)
	}

	// White tries to play at (4,4) which would be suicide
	// This should fail
	suicideMove := game.NewMove(4, 4, eng.White)
	err := game.MakeMove(suicideMove)
	if err == nil {
		t.Error("Suicide move should have been rejected")
	}
}

// TestKoRule tests basic Ko situation
func TestKoRule(t *testing.T) {
	game := engine.NewGame(9)

	// Create a Ko situation
	// Black stones
	if err := game.MakeMove(game.NewMove(4, 4, eng.Black)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(2, 2, eng.White)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(5, 5, eng.Black)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(4, 5, eng.White)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(6, 4, eng.Black)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(5, 3, eng.White)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(5, 4, eng.Black)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if err := game.MakeMove(game.NewMove(6, 5, eng.White)); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Black captures white stone at (5,4)
	if err := game.MakeMove(game.NewMove(4, 3, eng.Black)); err != nil {
		t.Fatalf("Capture move failed: %v", err)
	}

	// White immediately tries to recapture - should be blocked by Ko rule
	koMove := game.NewMove(5, 4, eng.White)
	err := game.MakeMove(koMove)
	if err == nil {
		t.Error("Ko rule should have prevented immediate recapture")
	}
}

// TestUndoRedo tests undo and redo functionality
func TestUndoRedo(t *testing.T) {
	game := engine.NewGame(9)

	// Make some moves
	moves := []struct{ x, y int }{
		{4, 4}, // black
		{5, 5}, // white
		{6, 6}, // black
	}

	for i, m := range moves {
		color := eng.Black
		if i%2 == 1 {
			color = eng.White
		}
		if err := game.MakeMove(game.NewMove(m.x, m.y, color)); err != nil {
			t.Fatalf("Move %d failed: %v", i+1, err)
		}
	}

	if game.MoveCount() != 3 {
		t.Errorf("Expected 3 moves, got %d", game.MoveCount())
	}

	// Test undo
	if !game.CanUndo() {
		t.Error("Should be able to undo")
	}

	if err := game.Undo(); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if game.MoveCount() != 2 {
		t.Errorf("After undo, expected 2 moves, got %d", game.MoveCount())
	}

	if game.CurrentTurn() != eng.Black {
		t.Errorf("After undo, expected Black's turn, got %v", game.CurrentTurn())
	}

	// Test redo
	if !game.CanRedo() {
		t.Error("Should be able to redo")
	}

	if err := game.Redo(); err != nil {
		t.Fatalf("Redo failed: %v", err)
	}

	if game.MoveCount() != 3 {
		t.Errorf("After redo, expected 3 moves, got %d", game.MoveCount())
	}

	// Test multiple undos
	game.Undo()
	game.Undo()
	game.Undo()

	if game.MoveCount() != 0 {
		t.Errorf("After 3 undos, expected 0 moves, got %d", game.MoveCount())
	}

	// Can't undo at start
	err := game.Undo()
	if err == nil {
		t.Error("Should not be able to undo at game start")
	}
}

// TestPassMechanism tests pass functionality
func TestPassMechanism(t *testing.T) {
	game := engine.NewGame(9)

	// Make a move
	if err := game.MakeMove(game.NewMove(4, 4, eng.Black)); err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	// White passes
	if err := game.Pass(); err != nil {
		t.Fatalf("Pass failed: %v", err)
	}

	if game.CurrentTurn() != eng.Black {
		t.Error("After white passes, it should be black's turn")
	}

	if game.IsGameOver() {
		t.Error("Game should not be over after one pass")
	}

	// Black passes - game should end
	if err := game.Pass(); err != nil {
		t.Fatalf("Pass failed: %v", err)
	}

	if !game.IsGameOver() {
		t.Error("Game should be over after both players pass")
	}
}

// TestResign tests resign functionality
func TestResign(t *testing.T) {
	game := engine.NewGame(9)

	// Make some moves
	game.MakeMove(game.NewMove(4, 4, eng.Black))
	game.MakeMove(game.NewMove(5, 5, eng.White))

	if game.IsGameOver() {
		t.Error("Game should not be over yet")
	}

	// Black resigns
	game.Resign()

	if !game.IsGameOver() {
		t.Error("Game should be over after resign")
	}

	// Can't make moves after resign
	err := game.MakeMove(game.NewMove(6, 6, eng.Black))
	if err == nil {
		t.Error("Should not be able to make moves after game is over")
	}
}

// TestScoring tests basic scoring functionality
func TestScoring(t *testing.T) {
	game := engine.NewGame(9)

	// Create a simple position with clear territories
	// Black stones
	game.MakeMove(game.NewMove(2, 2, eng.Black))
	game.MakeMove(game.NewMove(7, 7, eng.White))
	game.MakeMove(game.NewMove(2, 3, eng.Black))
	game.MakeMove(game.NewMove(7, 8, eng.White))
	game.MakeMove(game.NewMove(3, 2, eng.Black))
	game.MakeMove(game.NewMove(8, 7, eng.White))
	game.MakeMove(game.NewMove(3, 3, eng.Black))
	game.MakeMove(game.NewMove(8, 8, eng.White))

	score := game.GetScore()

	if score.BlackStones < 4 {
		t.Errorf("Black should have at least 4 stones, got %d", score.BlackStones)
	}

	if score.WhiteStones < 4 {
		t.Errorf("White should have at least 4 stones, got %d", score.WhiteStones)
	}

	// Test score with komi
	blackScore, whiteScore, winner := game.GetScoreWithKomi(6.5)

	if blackScore < 0 {
		t.Errorf("Black score should be non-negative, got %f", blackScore)
	}

	if whiteScore < 0 {
		t.Errorf("White score should be non-negative, got %f", whiteScore)
	}

	if winner != eng.Black && winner != eng.White {
		t.Errorf("Winner should be Black or White, got %v", winner)
	}
}

// TestTurnEnforcement tests that players must play in order
func TestTurnEnforcement(t *testing.T) {
	game := engine.NewGame(9)

	// Black's turn - white tries to play
	wrongTurnMove := game.NewMove(4, 4, eng.White)
	err := game.MakeMove(wrongTurnMove)
	if err == nil {
		t.Error("Should not allow white to play on black's turn")
	}

	// Black plays correctly
	if err := game.MakeMove(game.NewMove(4, 4, eng.Black)); err != nil {
		t.Fatalf("Valid move failed: %v", err)
	}

	// White's turn - black tries to play again
	wrongTurnMove2 := game.NewMove(5, 5, eng.Black)
	err = game.MakeMove(wrongTurnMove2)
	if err == nil {
		t.Error("Should not allow black to play on white's turn")
	}
}

// TestHistoryTruncation tests that making a move after undo truncates future history
func TestHistoryTruncation(t *testing.T) {
	game := engine.NewGame(9)

	// Make 3 moves
	game.MakeMove(game.NewMove(4, 4, eng.Black))
	game.MakeMove(game.NewMove(5, 5, eng.White))
	game.MakeMove(game.NewMove(6, 6, eng.Black))

	// Undo twice
	game.Undo()
	game.Undo()

	if game.MoveCount() != 1 {
		t.Errorf("After 2 undos, expected 1 move, got %d", game.MoveCount())
	}

	// Make a different move
	game.MakeMove(game.NewMove(3, 3, eng.White))

	// Should not be able to redo now
	if game.CanRedo() {
		t.Error("Should not be able to redo after making new move")
	}

	if game.MoveCount() != 2 {
		t.Errorf("Expected 2 moves after new branch, got %d", game.MoveCount())
	}
}

// TestComplexCapture tests multiple stone capture scenarios
func TestComplexCapture(t *testing.T) {
	game := engine.NewGame(9)

	// Test capturing multiple separate stones in one game
	// First capture: White stone at (5,4)
	game.MakeMove(game.NewMove(4, 4, eng.Black))
	game.MakeMove(game.NewMove(5, 4, eng.White)) // White stone to be captured
	game.MakeMove(game.NewMove(5, 3, eng.Black))
	game.MakeMove(game.NewMove(2, 2, eng.White)) // dummy elsewhere
	game.MakeMove(game.NewMove(5, 5, eng.Black))
	game.MakeMove(game.NewMove(3, 2, eng.White)) // dummy elsewhere
	game.MakeMove(game.NewMove(6, 4, eng.Black)) // Captures white at (5,4)

	// Verify first capture
	board := game.CurrentBoard()
	if board.At(5, 4) != eng.Empty {
		t.Error("White stone at (5,4) should have been captured")
	}

	// Continue game and test another capture at different location
	// Second capture: White stone at (8,8)
	game.MakeMove(game.NewMove(8, 8, eng.White)) // White stone to be captured
	game.MakeMove(game.NewMove(8, 7, eng.Black))
	game.MakeMove(game.NewMove(4, 2, eng.White)) // dummy
	game.MakeMove(game.NewMove(7, 8, eng.Black))
	game.MakeMove(game.NewMove(5, 2, eng.White)) // dummy
	game.MakeMove(game.NewMove(8, 9, eng.Black))
	game.MakeMove(game.NewMove(6, 2, eng.White)) // dummy
	game.MakeMove(game.NewMove(9, 8, eng.Black)) // Captures white at (8,8)

	// Verify second capture
	board = game.CurrentBoard()
	if board.At(8, 8) != eng.Empty {
		t.Error("White stone at (8,8) should have been captured")
	}

	// Verify we've made multiple moves and captures work correctly
	if game.MoveCount() != 15 {
		t.Errorf("Expected 15 moves, got %d", game.MoveCount())
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("Play at same position twice", func(t *testing.T) {
		game := engine.NewGame(9)

		move := game.NewMove(4, 4, eng.Black)
		if err := game.MakeMove(move); err != nil {
			t.Fatalf("First move failed: %v", err)
		}

		// Try to play at same position (after some other moves)
		game.MakeMove(game.NewMove(5, 5, eng.White))

		sameMove := game.NewMove(4, 4, eng.Black)
		err := game.MakeMove(sameMove)
		if err == nil {
			t.Error("Should not allow playing at occupied position")
		}
	})

	t.Run("Pass and undo", func(t *testing.T) {
		game := engine.NewGame(9)

		game.MakeMove(game.NewMove(4, 4, eng.Black))
		game.Pass() // white passes

		if game.CurrentTurn() != eng.Black {
			t.Error("Should be black's turn after white passes")
		}

		game.Undo() // undo the pass

		if game.CurrentTurn() != eng.White {
			t.Error("After undoing pass, should be white's turn again")
		}
	})

	t.Run("Multiple redos at end", func(t *testing.T) {
		game := engine.NewGame(9)

		game.MakeMove(game.NewMove(4, 4, eng.Black))

		// Try to redo when at end of history
		err := game.Redo()
		if err == nil {
			t.Error("Should not be able to redo at end of history")
		}
	})
}

// TestLongGame simulates a longer game sequence
func TestLongGame(t *testing.T) {
	game := engine.NewGame(9)

	// Play 20 moves
	for i := 1; i <= 10; i++ {
		x := 1 + (i % 8)
		y := 1 + ((i / 2) % 8)

		blackMove := game.NewMove(x, y, eng.Black)
		if err := game.MakeMove(blackMove); err != nil {
			t.Fatalf("Black move %d failed: %v", i, err)
		}

		x = 1 + ((i + 3) % 8)
		y = 1 + ((i/2 + 3) % 8)

		whiteMove := game.NewMove(x, y, eng.White)
		if err := game.MakeMove(whiteMove); err != nil {
			t.Fatalf("White move %d failed: %v", i, err)
		}
	}

	if game.MoveCount() != 20 {
		t.Errorf("Expected 20 moves, got %d", game.MoveCount())
	}

	// Test undo/redo in middle
	for i := 0; i < 5; i++ {
		game.Undo()
	}

	if game.MoveCount() != 15 {
		t.Errorf("After 5 undos, expected 15 moves, got %d", game.MoveCount())
	}

	for i := 0; i < 3; i++ {
		game.Redo()
	}

	if game.MoveCount() != 18 {
		t.Errorf("After 3 redos, expected 18 moves, got %d", game.MoveCount())
	}
}

// TestConcurrentUndoRedo tests undo/redo consistency
func TestConcurrentUndoRedo(t *testing.T) {
	game := engine.NewGame(9)

	// Build up history
	for i := 1; i <= 5; i++ {
		game.MakeMove(game.NewMove(i, i, eng.Black))
		game.MakeMove(game.NewMove(i+1, i, eng.White))
	}

	originalCount := game.MoveCount()

	// Undo all and redo all
	for i := 0; i < originalCount; i++ {
		if err := game.Undo(); err != nil {
			t.Fatalf("Undo %d failed: %v", i, err)
		}
	}

	if game.MoveCount() != 0 {
		t.Errorf("After undoing all, expected 0 moves, got %d", game.MoveCount())
	}

	for i := 0; i < originalCount; i++ {
		if err := game.Redo(); err != nil {
			t.Fatalf("Redo %d failed: %v", i, err)
		}
	}

	if game.MoveCount() != originalCount {
		t.Errorf("After redoing all, expected %d moves, got %d", originalCount, game.MoveCount())
	}
}

// TestGameStateConsistency tests that game state remains consistent
func TestGameStateConsistency(t *testing.T) {
	game := engine.NewGame(9)

	// Make moves and verify consistency
	moves := []struct {
		x, y  int
		color eng.Color
	}{
		{4, 4, eng.Black},
		{5, 5, eng.White},
		{6, 6, eng.Black},
	}

	for i, m := range moves {
		prevBoard := game.CurrentBoard()

		if err := game.MakeMove(game.NewMove(m.x, m.y, m.color)); err != nil {
			t.Fatalf("Move %d failed: %v", i+1, err)
		}

		newBoard := game.CurrentBoard()

		// Boards should be different objects
		if prevBoard == newBoard {
			t.Error("Board should be a new object after move")
		}

		// New board should have the move
		if newBoard.At(m.x, m.y) != m.color {
			t.Errorf("Move %d: expected color %v at (%d,%d), got %v",
				i+1, m.color, m.x, m.y, newBoard.At(m.x, m.y))
		}
	}
}

// TestBoardBoundaries tests moves at board edges
func TestBoardBoundaries(t *testing.T) {
	game := engine.NewGame(9)

	// Test corners and edges
	corners := []struct{ x, y int }{
		{1, 1}, // top-left
		{9, 1}, // top-right
		{1, 9}, // bottom-left
		{9, 9}, // bottom-right
	}

	for i, corner := range corners {
		color := eng.Black
		if i%2 == 1 {
			color = eng.White
		}
		if err := game.MakeMove(game.NewMove(corner.x, corner.y, color)); err != nil {
			t.Errorf("Corner move at (%d,%d) failed: %v", corner.x, corner.y, err)
		}
	}

	if game.MoveCount() != 4 {
		t.Errorf("Expected 4 moves, got %d", game.MoveCount())
	}
}

// TestMultipleGames tests that multiple games are independent
func TestMultipleGames(t *testing.T) {
	game1 := engine.NewGame(9)
	game2 := engine.NewGame(13)

	// Make moves in game1
	game1.MakeMove(game1.NewMove(4, 4, eng.Black))
	game1.MakeMove(game1.NewMove(5, 5, eng.White))

	// Make different moves in game2
	game2.MakeMove(game2.NewMove(7, 7, eng.Black))

	// Verify games are independent
	if game1.MoveCount() != 2 {
		t.Errorf("Game1 should have 2 moves, got %d", game1.MoveCount())
	}

	if game2.MoveCount() != 1 {
		t.Errorf("Game2 should have 1 move, got %d", game2.MoveCount())
	}

	if game1.Size() != 9 {
		t.Errorf("Game1 should be 9x9, got %d", game1.Size())
	}

	if game2.Size() != 13 {
		t.Errorf("Game2 should be 13x13, got %d", game2.Size())
	}
}

// TestScoreCalculation tests scoring in various scenarios
func TestScoreCalculation(t *testing.T) {
	t.Run("Empty board score", func(t *testing.T) {
		game := engine.NewGame(9)
		score := game.GetScore()

		if score.BlackStones != 0 {
			t.Errorf("Empty board should have 0 black stones, got %d", score.BlackStones)
		}

		if score.WhiteStones != 0 {
			t.Errorf("Empty board should have 0 white stones, got %d", score.WhiteStones)
		}
	})

	t.Run("Score after captures", func(t *testing.T) {
		game := engine.NewGame(9)

		// Set up a capture scenario
		game.MakeMove(game.NewMove(4, 4, eng.Black))
		game.MakeMove(game.NewMove(5, 4, eng.White))
		game.MakeMove(game.NewMove(5, 3, eng.Black))
		game.MakeMove(game.NewMove(2, 2, eng.White))
		game.MakeMove(game.NewMove(5, 5, eng.Black))
		game.MakeMove(game.NewMove(3, 2, eng.White))
		game.MakeMove(game.NewMove(6, 4, eng.Black)) // Captures white

		score := game.GetScore()

		// After capture, white should have fewer stones
		if score.WhiteStones != 2 {
			t.Errorf("Expected 2 white stones after capture, got %d", score.WhiteStones)
		}
	})
}

// TestUndoAfterCapture tests that undo correctly restores captured stones
func TestUndoAfterCapture(t *testing.T) {
	game := engine.NewGame(9)

	// Set up a capture
	game.MakeMove(game.NewMove(4, 4, eng.Black))
	game.MakeMove(game.NewMove(5, 4, eng.White))
	game.MakeMove(game.NewMove(5, 3, eng.Black))
	game.MakeMove(game.NewMove(2, 2, eng.White))
	game.MakeMove(game.NewMove(5, 5, eng.Black))
	game.MakeMove(game.NewMove(3, 2, eng.White))
	game.MakeMove(game.NewMove(6, 4, eng.Black)) // Captures white at (5,4)

	// Verify capture
	board := game.CurrentBoard()
	if board.At(5, 4) != eng.Empty {
		t.Error("White stone should have been captured")
	}

	// Undo the capture
	game.Undo()

	// Verify white stone is restored
	board = game.CurrentBoard()
	if board.At(5, 4) != eng.White {
		t.Error("White stone should be restored after undo")
	}
}

// TestPassResets tests that pass flags are reset after a move
func TestPassResets(t *testing.T) {
	game := engine.NewGame(9)

	// Black plays
	game.MakeMove(game.NewMove(4, 4, eng.Black))

	// White passes
	game.Pass()

	// Black passes (game should end)
	game.Pass()

	if !game.IsGameOver() {
		t.Error("Game should be over after both pass")
	}

	// Undo both passes
	game.Undo()
	game.Undo()

	// Black makes a move instead
	if err := game.MakeMove(game.NewMove(5, 5, eng.Black)); err != nil {
		t.Fatalf("Move after undo failed: %v", err)
	}

	// Game should not be over
	if game.IsGameOver() {
		t.Error("Game should not be over after move following undo of passes")
	}

	// White passes
	game.Pass()

	// Black makes a move (not pass)
	game.MakeMove(game.NewMove(6, 6, eng.Black))

	// White passes again
	game.Pass()

	// Game should not be over (black didn't pass consecutively)
	if game.IsGameOver() {
		t.Error("Game should not be over - passes were not consecutive")
	}
}
