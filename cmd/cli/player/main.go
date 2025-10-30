package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/awesohame/gogo/pkg/engine"
)

func main() {
	fmt.Println("GoGo CLI 2-player dev tool (9x9)")
	reader := bufio.NewReader(os.Stdin)
	game := engine.NewGame(9)
	for {
		fmt.Println("\nCurrent board:")
		fmt.Print(game.CurrentBoard().String())
		if game.IsGameOver() {
			fmt.Println("Game over!")
			score := game.GetScore()
			fmt.Printf("Score - Black: %d, White: %d\n", score.Black, score.White)
			break
		}
		turn := game.CurrentTurn()
		turnStr := "Black"
		if turn == 2 { // engine.White
			turnStr = "White"
		}
		fmt.Printf("%s's turn. Enter move (x y), or command (pass, resign, undo, redo, exit): ", turnStr)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch strings.ToLower(input) {
		case "pass":
			err := game.Pass()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "resign":
			game.Resign()
			fmt.Printf("%s resigned. Game over!\n", turnStr)
		case "undo":
			err := game.Undo()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "redo":
			err := game.Redo()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "exit":
			fmt.Println("Exiting game.")
			return
		default:
			parts := strings.Fields(input)
			if len(parts) == 2 {
				x, err1 := strconv.Atoi(parts[0])
				y, err2 := strconv.Atoi(parts[1])
				if err1 != nil || err2 != nil || x < 1 || y < 1 || x > game.Size() || y > game.Size() {
					fmt.Println("Invalid coordinates. Enter x y with 1 <= x,y <=", game.Size())
					continue
				}
				move := game.NewMove(x, y, turn)
				err := game.MakeMove(move)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command. Enter move as 'x y' or a valid command.")
			}
		}
	}
}
