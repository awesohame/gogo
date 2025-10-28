package main

import (
	"fmt"

	"github.com/awesohame/gogo/internal/engine"
)

func main() {
	fmt.Println("--- GoGo Engine Dev Test ---")

	board := engine.NewBoard(9)
	fmt.Println("1. Initial 9x9 Board:")
	fmt.Println(board)

	var err error

	fmt.Println("\n2. Placing Black at (5, 5)...")
	move1 := engine.Move{Point: board.ToPoint(5, 5), Color: engine.Black}
	board, err = board.ApplyMove(move1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n3. Placing White at (5, 6)...")
	move2 := engine.Move{Point: board.ToPoint(5, 6), Color: engine.White}
	board, err = board.ApplyMove(move2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n4. Placing Black at (4, 5) to test friendly merge...")
	move3 := engine.Move{Point: board.ToPoint(4, 5), Color: engine.Black}
	board, err = board.ApplyMove(move3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n5. Placing Black at (6, 6) to surround White...")
	move4 := engine.Move{Point: board.ToPoint(6, 6), Color: engine.Black}
	board, err = board.ApplyMove(move4)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n5. Placing Black at (4, 6) to surround White...")
	move5 := engine.Move{Point: board.ToPoint(4, 6), Color: engine.Black}
	board, err = board.ApplyMove(move5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n5. Placing Black at (3, 6) to surround White...")
	move6 := engine.Move{Point: board.ToPoint(3, 6), Color: engine.Black}
	board, err = board.ApplyMove(move6)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n5. Placing Black at (4, 7) to surround White...")
	move7 := engine.Move{Point: board.ToPoint(4, 7), Color: engine.Black}
	board, err = board.ApplyMove(move7)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n6. Placing Black at (5, 7) to capture White...")
	lastmove := engine.Move{Point: board.ToPoint(5, 7), Color: engine.Black}
	board, err = board.ApplyMove(lastmove)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	// Test suicide detection
	fmt.Println("\n7. Testing suicide detection...")
	fmt.Println("Attempting to place White at (5, 6) (where it was captured)...")
	suicideMove := engine.Move{Point: board.ToPoint(5, 6), Color: engine.White}
	_, err = board.ApplyMove(suicideMove)
	if err != nil {
		fmt.Printf("Correctly rejected: %v\n", err)
	} else {
		fmt.Println("ERROR: Suicidal move was allowed!")
	}

	// Test Ko detection
	fmt.Println("\n8. Testing Ko detection...")
	// Create a simple Ko scenario on a fresh board
	koBoard := engine.NewBoard(9)

	// Set up a Ko position
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(3, 3), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 3), Color: engine.White})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 2), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(3, 4), Color: engine.White})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(5, 3), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(5, 4), Color: engine.White})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(8, 8), Color: engine.Black})
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 5), Color: engine.White})

	fmt.Println("Ko setup board:")
	fmt.Println(koBoard)

	// Black captures at (4,4) - this should create a Ko
	fmt.Println("Black captures Black at (4, 4)...")
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 4), Color: engine.Black})

	// Now White tries to immediately recapture - this should be rejected as Ko
	fmt.Println("White attempts immediate recapture (Ko violation)...")
	_, err = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 3), Color: engine.White})
	if err != nil {
		fmt.Printf("Correctly rejected Ko: %v\n", err)
	} else {
		fmt.Println("ERROR: Ko violation was allowed!")
	}

	// make another black move elsewhere
	koBoard, _ = koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(8, 7), Color: engine.Black})
	fmt.Println(koBoard)
	fmt.Println("White tries again after another move, no violation expected...")
	finalBoard, err := koBoard.ApplyMove(engine.Move{Point: koBoard.ToPoint(4, 3), Color: engine.White})
	if err == nil {
		fmt.Println("No Ko violation, move allowed.")
		koBoard = finalBoard
		fmt.Println("Board after White recapture:")
		fmt.Println(koBoard)
	} else {
		fmt.Printf("Incorrectly rejected Ko: %v\n", err)
	}

	fmt.Println("\n--- Test Complete ---")
}
