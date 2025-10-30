package ai

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/awesohame/gogo/internal/engine"
)

// MCTSBot implements Monte Carlo Tree Search
type MCTSBot struct {
	MaxSimulations int     // no. of simulations to run
	TimeLimit      float64 // time limit in sec (if 0, uses MaxSimulations)
	ExplorationC   float64 // UCB exploration const (sqrt 2)
}

// NewMCTSBot creates a new MCTS bot
func NewMCTSBot(simulations int) *MCTSBot {
	return &MCTSBot{
		MaxSimulations: simulations,
		TimeLimit:      0, // no time lim by default
		ExplorationC:   math.Sqrt(2),
	}
}

// implements the Bot interface using MCTS
func (bot *MCTSBot) SelectMove(board *engine.Board, color engine.Color) engine.Move {
	previousColor := opponentColor(color)
	root := newMCTSNode(nil, engine.Move{Point: -1, Color: previousColor}, board, previousColor)

	startTime := time.Now()
	simulations := 0

	// run simulations until we hit the lim
	for {
		if bot.TimeLimit > 0 {
			if time.Since(startTime).Seconds() >= bot.TimeLimit {
				break
			}
		} else if simulations >= bot.MaxSimulations {
			break
		}

		// MCTS -> selection, expansion, simulation, backpropagation
		node := root.selectNode(bot.ExplorationC)
		winner := node.simulate()
		node.backpropagate(winner)

		simulations++
	}

	// choose the best move based on visit count
	bestChild := root.bestChild(0) // 0 means no exploration
	if bestChild == nil {
		// pass when no legal moves
		return engine.Move{Point: -1, Color: color}
	}

	fmt.Printf("MCTS: %d simulations, selected move with %d visits (%.1f%% win rate)\n",
		simulations, bestChild.visits, 100.0*bestChild.wins/float64(bestChild.visits))

	return bestChild.move
}

// MCTSNode represents a node in the MCTS tree
type MCTSNode struct {
	parent   *MCTSNode
	children []*MCTSNode
	move     engine.Move
	board    *engine.Board
	color    engine.Color // color of the player who just moved to reach this state

	visits       int
	wins         float64 // wins from the perspective of the parent's color
	untriedMoves []engine.Move
}

// newMCTSNode creates a new MCTS node
func newMCTSNode(parent *MCTSNode, move engine.Move, board *engine.Board, color engine.Color) *MCTSNode {
	node := &MCTSNode{
		parent:   parent,
		children: make([]*MCTSNode, 0),
		move:     move,
		board:    board,
		color:    color,
		visits:   0,
		wins:     0,
	}

	// get all legal moves for the next player
	nextColor := opponentColor(color)
	node.untriedMoves = getLegalMoves(board, nextColor)

	return node
}

// traverses the tree using UCB1 until a leaf node
func (n *MCTSNode) selectNode(explorationC float64) *MCTSNode {
	current := n

	for len(current.untriedMoves) == 0 && len(current.children) > 0 {
		current = current.bestChild(explorationC)
	}

	// if there are untried moves, expand
	if len(current.untriedMoves) > 0 {
		return current.expand()
	}

	return current
}

// adds a new child node for an untried move
func (n *MCTSNode) expand() *MCTSNode {
	// pick random untried move
	idx := rand.Intn(len(n.untriedMoves))
	move := n.untriedMoves[idx]

	// remove from untried moves
	n.untriedMoves = append(n.untriedMoves[:idx], n.untriedMoves[idx+1:]...)

	// apply move
	newBoard, err := n.board.ApplyMove(move)
	if err != nil {
		// fix getLegalMoves if this happens
		return n
	}

	// create child node
	childNode := newMCTSNode(n, move, newBoard, move.Color)
	n.children = append(n.children, childNode)

	return childNode
}

// runs a random playout from this node and returns the winner
func (n *MCTSNode) simulate() engine.Color {
	board := n.board
	currentColor := opponentColor(n.color)
	passCount := 0
	maxMoves := 200 // limit simulation length
	moveCount := 0

	for moveCount < maxMoves {
		// fast legal move generation with caching
		legalMoves := getFastLegalMoves(board, currentColor, 30) // limit candidates

		if len(legalMoves) == 0 {
			passCount++
			if passCount >= 2 {
				// both players passed, game over
				break
			}
			currentColor = opponentColor(currentColor)
			continue
		}

		passCount = 0
		moveCount++

		// pick a random move
		move := legalMoves[rand.Intn(len(legalMoves))]
		newBoard, err := board.ApplyMove(move)
		if err != nil {
			//skip invalid move
			continue
		}

		board = newBoard
		currentColor = opponentColor(currentColor)
	}

	// get winner by score
	blackScore, whiteScore, winner := board.CalculateScoreWithKomi(6.5)
	_ = blackScore
	_ = whiteScore

	return winner
}

// updates statistics up the tree
func (n *MCTSNode) backpropagate(winner engine.Color) {
	current := n

	for current != nil {
		current.visits++

		// Update wins from parent's perspective
		if current.parent != nil {
			parentColor := current.parent.color
			switch winner {
			case parentColor:
				current.wins += 1.0
			case engine.Empty:
				current.wins += 0.5 // draw
			}
		}

		current = current.parent
	}
}

// returns child with highest UCB1 score
func (n *MCTSNode) bestChild(explorationC float64) *MCTSNode {
	if len(n.children) == 0 {
		return nil
	}

	var bestChild *MCTSNode
	bestScore := -math.MaxFloat64

	for _, child := range n.children {
		score := child.ucb1Score(explorationC)
		if score > bestScore {
			bestScore = score
			bestChild = child
		}
	}

	return bestChild
}

// calcs the UCB1 score for this node
func (n *MCTSNode) ucb1Score(explorationC float64) float64 {
	if n.visits == 0 {
		return math.Inf(1) // unvisited nodes first
	}

	exploitation := n.wins / float64(n.visits)
	exploration := explorationC * math.Sqrt(math.Log(float64(n.parent.visits))/float64(n.visits))

	return exploitation + exploration
}

// returns all legal moves for the given color
func getLegalMoves(board *engine.Board, color engine.Color) []engine.Move {
	moves := make([]engine.Move, 0, 30) // preallocate capacity
	size := board.Size()

	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			point := board.ToPoint(x, y)

			// only consider empty points
			if board.At(x, y) != engine.Empty {
				continue
			}

			move := engine.Move{Point: point, Color: color}

			// try to apply move
			_, err := board.ApplyMove(move)
			if err == nil {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// returns a limited set of legal moves for fast simulation
func getFastLegalMoves(board *engine.Board, color engine.Color, maxMoves int) []engine.Move {
	moves := make([]engine.Move, 0, maxMoves)
	size := board.Size()

	// random sampling for speed
	tried := make(map[engine.Point]bool)
	attempts := 0
	maxAttempts := size * size

	for len(moves) < maxMoves && attempts < maxAttempts {
		// random point
		x := rand.Intn(size) + 1
		y := rand.Intn(size) + 1
		point := board.ToPoint(x, y)

		if tried[point] {
			attempts++
			continue
		}
		tried[point] = true
		attempts++

		// only consider empty points
		if board.At(x, y) != engine.Empty {
			continue
		}

		// dont fill own eyes
		if IsEyeFillingMove(board, engine.Move{Point: point, Color: color}) {
			continue
		}

		move := engine.Move{Point: point, Color: color}

		// try to apply move
		_, err := board.ApplyMove(move)
		if err == nil {
			moves = append(moves, move)
		}
	}

	// if we didnt find enough moves, use search to get moves
	if len(moves) < 5 {
		for y := 1; y <= size && len(moves) < maxMoves; y++ {
			for x := 1; x <= size && len(moves) < maxMoves; x++ {
				point := board.ToPoint(x, y)

				if tried[point] || board.At(x, y) != engine.Empty {
					continue
				}

				move := engine.Move{Point: point, Color: color}
				_, err := board.ApplyMove(move)
				if err == nil {
					moves = append(moves, move)
				}
			}
		}
	}

	return moves
}
