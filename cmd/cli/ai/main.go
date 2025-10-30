package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/awesohame/gogo/internal/ai"
	"github.com/awesohame/gogo/internal/engine"
)

func main() {
	fmt.Println("=== Go Game: Play against MCTS AI ===")
	fmt.Println()

	// game setup
	fmt.Print("Board size (9, 13, or 19): ")
	size := readInt(9)

	fmt.Print("AI strength - number of simulations (100-5000): ")
	simulations := readInt(500)

	fmt.Print("Play as (1=Black, 2=White): ")
	humanColor := engine.Black
	if readInt(1) == 2 {
		humanColor = engine.White
	}

	// create game
	board := engine.NewBoard(size)
	bot := ai.NewMCTSBot(simulations)

	currentColor := engine.Black
	passCount := 0
	moveNumber := 1

	fmt.Println("\n=== Game Start ===")
	fmt.Println("Commands: <x> <y> to play, 'pass' to pass, 'quit' to exit")
	fmt.Println()

	// game loop
	for {
		fmt.Printf("\n--- Move %d: %v to play ---\n", moveNumber, colorName(currentColor))
		fmt.Println(board.String())

		var move engine.Move

		if currentColor == humanColor {
			// human
			fmt.Print("Your move: ")
			move = readMove(board, currentColor)

			if move.Point == -2 {
				// quit
				fmt.Println("Game ended by player.")
				return
			}
		} else {
			// AI
			fmt.Println("AI is thinking...")
			move = bot.SelectMove(board, currentColor)
		}

		// process move
		if move.Point == -1 {
			// pass
			fmt.Printf("%v passes\n", colorName(currentColor))
			passCount++

			if passCount >= 2 {
				break // game over
			}

			currentColor = opponentColor(currentColor)
			moveNumber++
			continue
		}

		passCount = 0

		// apply move
		newBoard, err := board.ApplyMove(move)
		if err != nil {
			fmt.Printf("Invalid move: %v\n", err)
			if currentColor != humanColor {
				// AI made invalid move (fix getLegalMoves)
				fmt.Println("AI error - ending game")
				break
			}
			continue
		}

		x, y := board.ToXY(move.Point)
		fmt.Printf("%v plays at (%d, %d)\n", colorName(currentColor), x, y)

		board = newBoard
		currentColor = opponentColor(currentColor)
		moveNumber++
	}

	// game over
	fmt.Println("\n=== Game Over ===")
	fmt.Println("Final board:")
	fmt.Println(board.String())

	blackScore, whiteScore, winner := board.CalculateScoreWithKomi(6.5)
	fmt.Printf("\nFinal Score:\n")
	fmt.Printf("Black: %.1f\n", blackScore)
	fmt.Printf("White: %.1f (with 6.5 komi)\n", whiteScore)
	fmt.Printf("\nWinner: %v\n", colorName(winner))

	if winner == humanColor {
		fmt.Println("Congratulations! You won!")
	} else if winner == opponentColor(humanColor) {
		fmt.Println("AI wins! Better luck next time!")
	} else {
		fmt.Println("It's a draw!")
	}
}

func readInt(defaultVal int) int {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if line == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(line)
	if err != nil {
		return defaultVal
	}

	return val
}

func readMove(board *engine.Board, color engine.Color) engine.Move {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))

		if line == "quit" || line == "exit" {
			return engine.Move{Point: -2, Color: color}
		}

		if line == "pass" {
			return engine.Move{Point: -1, Color: color}
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Print("Invalid format. Try '<x> <y>' or 'pass': ")
			continue
		}

		x, err1 := strconv.Atoi(parts[0])
		y, err2 := strconv.Atoi(parts[1])

		if err1 != nil || err2 != nil {
			fmt.Print("Invalid coordinates. Try again: ")
			continue
		}

		if x < 1 || x > board.Size() || y < 1 || y > board.Size() {
			fmt.Printf("Coordinates must be between 1 and %d. Try again: ", board.Size())
			continue
		}

		point := board.ToPoint(x, y)
		return engine.Move{Point: point, Color: color}
	}
}

func colorName(c engine.Color) string {
	switch c {
	case engine.Black:
		return "Black"
	case engine.White:
		return "White"
	default:
		return "Empty"
	}
}

func opponentColor(c engine.Color) engine.Color {
	if c == engine.Black {
		return engine.White
	}
	return engine.Black
}
