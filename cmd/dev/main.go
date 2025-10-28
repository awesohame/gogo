package main

import (
	"fmt"

	"github.com/awesohame/gogo/internal/engine"
)

func main() {
	fmt.Println("gogo engine quick test")
	testBasicRules()
	testScoring()
	deadStoneTest()
}

func testBasicRules() {
	board := engine.NewBoard(9)
	fmt.Println("init 9x9 board")
	fmt.Println(board)

	var err error

	fmt.Println("black at 5,5")
	move1 := engine.Move{Point: board.ToPoint(5, 5), Color: engine.Black}
	board, _ = board.ApplyMove(move1)
	fmt.Println(board)

	fmt.Println("white at 5,6")
	move2 := engine.Move{Point: board.ToPoint(5, 6), Color: engine.White}
	board, _ = board.ApplyMove(move2)
	fmt.Println(board)

	fmt.Println("black at 4,5 (merge)")
	move3 := engine.Move{Point: board.ToPoint(4, 5), Color: engine.Black}
	board, _ = board.ApplyMove(move3)
	fmt.Println(board)

	fmt.Println("black at 6,6 (surround)")
	move4 := engine.Move{Point: board.ToPoint(6, 6), Color: engine.Black}
	board, _ = board.ApplyMove(move4)
	fmt.Println(board)

	fmt.Println("black at 4,6")
	move5 := engine.Move{Point: board.ToPoint(4, 6), Color: engine.Black}
	board, _ = board.ApplyMove(move5)
	fmt.Println(board)

	fmt.Println("black at 3,6")
	move6 := engine.Move{Point: board.ToPoint(3, 6), Color: engine.Black}
	board, _ = board.ApplyMove(move6)
	fmt.Println(board)

	fmt.Println("black at 4,7")
	move7 := engine.Move{Point: board.ToPoint(4, 7), Color: engine.Black}
	board, _ = board.ApplyMove(move7)
	fmt.Println(board)

	fmt.Println("black at 5,7 (capture)")
	lastmove := engine.Move{Point: board.ToPoint(5, 7), Color: engine.Black}
	board, _ = board.ApplyMove(lastmove)
	fmt.Println(board)

	// suicide test
	fmt.Println("try white at 5,6 (should fail)")
	suicideMove := engine.Move{Point: board.ToPoint(5, 6), Color: engine.White}
	_, err = board.ApplyMove(suicideMove)
	if err != nil {
		fmt.Println("suicide correctly rejected")
	} else {
		fmt.Println("suicide move allowed (bug)")
	}

	// ko test
	fmt.Println("ko test setup")
	koBoard := engine.NewBoard(9)
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(3, 3), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 3), Color: engine.White})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 2), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(3, 4), Color: engine.White})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(5, 3), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(5, 4), Color: engine.White})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(8, 8), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 5), Color: engine.White})
	fmt.Println(koBoard)

	fmt.Println("black at 4,4 (ko capture)")
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 4), Color: engine.Black})

	fmt.Println("white tries ko recapture (should fail)")
	_, err = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 3), Color: engine.White})
	if err != nil {
		fmt.Println("ko correctly rejected")
	} else {
		fmt.Println("ko move allowed (bug)")
	}

	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(8, 7), Color: engine.Black})
	fmt.Println(koBoard)
	fmt.Println("white tries again at 4,3 (should work)")
	finalBoard, err := koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 3), Color: engine.White})
	if err == nil {
		fmt.Println("ko move allowed after other move")
		koBoard = finalBoard
		fmt.Println(koBoard)
	} else {
		fmt.Println("ko move rejected (bug)")
	}

	fmt.Println("done with basic rules test")
}

func testScoring() {
	fmt.Println("scoring test")
	board := engine.NewBoard(9)
	// black surrounds top left 4x4, leaves 2x2 empty
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(9, 9), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(2, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(8, 9), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(3, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(9, 8), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(4, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(8, 8), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(4, 2), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(9, 7), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(4, 3), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(8, 7), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(4, 4), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(7, 7), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(3, 4), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(7, 8), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(2, 4), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(7, 9), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 4), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 9), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 3), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(2, 9), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 2), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(3, 9), Color: engine.White})
	fmt.Println("final board:")
	fmt.Println(board)
	score := board.CalculateChineseScore()
	fmt.Println("black stones:", score.BlackStones)
	fmt.Println("black territory:", score.BlackArea)
	fmt.Println("black total:", score.Black)
	fmt.Println("white stones:", score.WhiteStones)
	fmt.Println("white territory:", score.WhiteArea)
	fmt.Println("white total:", score.White)
	fmt.Println("dame:", score.DamePoints)
	totalPoints := score.Black + score.White + score.DamePoints
	fmt.Println("check:", score.Black, "+", score.White, "+", score.DamePoints, "=", totalPoints, "should be", board.Size()*board.Size())
	if score.Black > score.White {
		fmt.Println("black wins by", score.Black-score.White)
	} else if score.White > score.Black {
		fmt.Println("white wins by", score.White-score.Black)
	} else {
		fmt.Println("draw")
	}
	fmt.Println("komi test (7.5)")
	blackScore, whiteScore, winner := board.CalculateScoreWithKomi(7.5)
	fmt.Println("black:", blackScore)
	fmt.Println("white:", whiteScore)
	switch winner {
	case engine.Black:
		fmt.Println("black wins by", blackScore-whiteScore)
	case engine.White:
		fmt.Println("white wins by", whiteScore-blackScore)
	default:
		fmt.Println("draw")
	}
}

func deadStoneTest() {
	fmt.Println("\nlife/death test")
	board := engine.NewBoard(9)

	// create a dead white group in black territory
	// black surrounds a white stone that's in atari
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(2, 2), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 2), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(3, 2), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(2, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(2, 3), Color: engine.Black})

	// add more black stones to secure territory
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(3, 1), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(1, 3), Color: engine.Black})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(3, 3), Color: engine.Black})

	// add a live white group elsewhere
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(7, 7), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(8, 7), Color: engine.White})
	board, _ = board.ApplyMove(engine.Move{Point: board.ToPoint(7, 8), Color: engine.White})

	fmt.Println("board with dead white stone at 2,2:")
	fmt.Println(board)

	score := board.CalculateChineseScore()
	fmt.Println("scoring (should remove dead white stone)")
	fmt.Println("black stones:", score.BlackStones)
	fmt.Println("black territory:", score.BlackArea)
	fmt.Println("black total:", score.Black)
	fmt.Println("white stones:", score.WhiteStones, "(dead stones removed)")
	fmt.Println("white territory:", score.WhiteArea)
	fmt.Println("white total:", score.White)
	fmt.Println("dame:", score.DamePoints)
}
