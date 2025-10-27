package main

import (
	"fmt"
	"log"

	"github.com/awesohame/gogo/internal/engine"
)

func main() {
	// create board
	board := engine.NewBoard(9)
	fmt.Println("--- Initial Empty Board ---")
	fmt.Println(board)

	// place black stone at (5,5)
	move := engine.Move{
		Point: board.ToPoint(5, 5),
		Color: engine.Black,
	}

	newBoard, err := board.ApplyMove(move)
	if err != nil {
		log.Fatalf("Error applying move: %v", err)
	}

	// print board
	fmt.Println("\n--- Board After Placing Black at (5,5) ---")
	fmt.Println(newBoard)
}
