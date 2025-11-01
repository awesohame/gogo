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
	MaxSimulations int       // no. of simulations to run
	TimeLimit      float64   // time limit in sec (if 0, uses MaxSimulations)
	ExplorationC   float64   // UCB exploration const (sqrt 2)
	ReuseTree      bool      // whether to reuse tree between moves
	lastRoot       *MCTSNode // root from previous move for tree reuse
}

// NewMCTSBot creates a new MCTS bot
func NewMCTSBot(simulations int) *MCTSBot {
	return &MCTSBot{
		MaxSimulations: simulations,
		TimeLimit:      0, // no time lim by default
		ExplorationC:   math.Sqrt(2),
		ReuseTree:      true, // tree reuse by default
	}
}

// implements the Bot interface using MCTS
func (bot *MCTSBot) SelectMove(board *engine.Board, color engine.Color) engine.Move {
	previousColor := opponentColor(color)

	// Try to reuse tree from previous move
	var root *MCTSNode
	if bot.ReuseTree && bot.lastRoot != nil {
		// Find child matching current board position
		root = bot.findMatchingChild(bot.lastRoot, board)
		if root != nil {
			root.parent = nil // detach from old tree
		}
	}

	// If no reuse, create new root
	if root == nil {
		root = newMCTSNode(nil, engine.Move{Point: -1, Color: previousColor}, board, previousColor)
	}

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

		// adaptive exploration: reduce exploration as we get more confident
		explorationC := bot.ExplorationC
		if simulations > bot.MaxSimulations/2 {
			explorationC *= 0.8 // reduce exploration in later phase
		}

		// MCTS -> selection, expansion, simulation, backpropagation
		node := root.selectNode(explorationC)
		winner := node.simulate()
		node.backpropagate(winner)

		simulations++
	}

	// choose the best move based on visit count
	bestChild := root.bestChild(0) // 0 means no exploration
	if bestChild == nil {
		// pass when no legal moves
		bot.lastRoot = nil
		return engine.Move{Point: -1, Color: color}
	}

	// Save tree for reuse
	if bot.ReuseTree {
		bot.lastRoot = bestChild
	}

	fmt.Printf("MCTS: %d simulations, selected move with %d visits (%.1f%% win rate)\n",
		simulations, bestChild.visits, 100.0*bestChild.wins/float64(bestChild.visits))

	return bestChild.move
}

// attempts to find a child node matching the current board
func (bot *MCTSBot) findMatchingChild(node *MCTSNode, board *engine.Board) *MCTSNode {
	// zobrist hash to find matching child position
	targetHash := board.Hash()
	for _, child := range node.children {
		if child.board.Hash() == targetHash {
			return child
		}
	}
	return nil
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

	// get all legal moves for the next player (lazily, only if needed)
	nextColor := opponentColor(color)
	node.untriedMoves = getLegalMovesFast(board, nextColor)

	return node
}

// faster legal move generation that doesn't validate every move
func getLegalMovesFast(board *engine.Board, color engine.Color) []engine.Move {
	moves := make([]engine.Move, 0, 40)
	size := board.Size()

	// prioritize center and corner/edge positions early
	priorityMoves := make([]engine.Move, 0, 20)
	normalMoves := make([]engine.Move, 0, 20)

	center := size / 2

	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			// only consider empty points
			if board.At(x, y) != engine.Empty {
				continue
			}

			point := board.ToPoint(x, y)
			move := engine.Move{Point: point, Color: color}

			// quick heuristic: skip obvious eye fills
			if IsEyeFillingMove(board, move) {
				continue
			}

			// prioritize center area and corners/edges
			distFromCenter := abs(x-center) + abs(y-center)
			isEdge := x == 1 || x == size || y == 1 || y == size

			if distFromCenter <= 2 || isEdge {
				priorityMoves = append(priorityMoves, move)
			} else {
				normalMoves = append(normalMoves, move)
			}
		}
	}

	// return priority moves first, then normal moves
	moves = append(moves, priorityMoves...)
	moves = append(moves, normalMoves...)

	return moves
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
	if len(n.untriedMoves) == 0 {
		return n
	}

	// pick random untried move
	idx := rand.Intn(len(n.untriedMoves))
	move := n.untriedMoves[idx]

	// remove from untried moves (swap with last for efficiency)
	n.untriedMoves[idx] = n.untriedMoves[len(n.untriedMoves)-1]
	n.untriedMoves = n.untriedMoves[:len(n.untriedMoves)-1]

	// apply move
	newBoard, err := n.board.ApplyMove(move)
	if err != nil {
		// move became illegal, try another if available
		if len(n.untriedMoves) > 0 {
			return n.expand()
		}
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
	maxMoves := 150
	moveCount := 0

	// early termination score threshold
	earlyCheckInterval := 30

	for moveCount < maxMoves {
		// early termination check - if one side is clearly winning, end simulation
		if moveCount > 0 && moveCount%earlyCheckInterval == 0 {
			blackScore, whiteScore, _ := board.CalculateScoreWithKomi(6.5)
			scoreDiff := blackScore - whiteScore
			// if score difference is huge, end early
			if scoreDiff > 20 || scoreDiff < -20 {
				break
			}
		}

		// fast legal move generation with caching
		legalMoves := getFastLegalMoves(board, currentColor, 25) // reduced from 30

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
	_, _, winner := board.CalculateScoreWithKomi(0)

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

// returns a limited set of legal moves for fast simulation
func getFastLegalMoves(board *engine.Board, color engine.Color, maxMoves int) []engine.Move {
	moves := make([]engine.Move, 0, maxMoves)
	size := board.Size()

	// use map for correct indexing with internal board points
	tried := make(map[engine.Point]bool, size*size)
	attempts := 0
	maxAttempts := min(size*size, maxMoves*3) // limit attempts

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

		// dont fill own eyes (skip check in late game for speed)
		if len(moves) < maxMoves/2 && IsEyeFillingMove(board, engine.Move{Point: point, Color: color}) {
			continue
		}

		move := engine.Move{Point: point, Color: color}

		// try to apply move
		_, err := board.ApplyMove(move)
		if err == nil {
			moves = append(moves, move)
		}
	}

	// if we didnt find enough moves, use systematic search
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
