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

	fmt.Println("\n6. Placing Black at (5, 7) to capture White...")
	move6 := engine.Move{Point: board.ToPoint(5, 7), Color: engine.Black}
	board, err = board.ApplyMove(move6)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(board)

	fmt.Println("\n--- Test Complete ---")
}
