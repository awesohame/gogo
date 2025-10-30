package engine

import (
	"errors"
	"strings"
)

// Color is the state of a point
type Color int8

const (
	Empty Color = iota
	Black
	White
	Border // Border for internal calc
)

// Point is a board coordinate
type Point int16

// Board represents game state
type Board struct {
	points       []Color
	size         int
	internalSize int
	koPoint      Point    // active Ko point
	koHash       uint64   // Zobrist hash for Ko detection
	history      []uint64 // Zobrist hash history for superko
	groups       map[Point]*Group
	dsu          *DSU
	nextGroupID  int
}

// NewBoard inits new board of a given size
func NewBoard(size int) *Board {
	if size <= 0 {
		size = 9 // default board size
	}

	internalSize := size + 2
	points := make([]Color, internalSize*internalSize)

	// init all points to Empty and set borders
	for i := range points {
		points[i] = Empty
	}

	// set top and bottom borders
	for i := 0; i < internalSize; i++ {
		points[i] = Border
		points[i+internalSize*(internalSize-1)] = Border
	}

	// set left and right borders
	for i := 0; i < internalSize; i++ {
		points[i*internalSize] = Border
		points[i*internalSize+internalSize-1] = Border
	}

	return &Board{
		points:       points,
		size:         size,
		internalSize: internalSize,
		groups:       make(map[Point]*Group),
		dsu:          NewDSU(size * size * 100), // extra room for MCTS
		history:      make([]uint64, 0),
		koPoint:      -1, // use -1 for no active Ko point
	}
}

// creates a deep copy of curr board state
func (b *Board) copy() *Board {
	newPoints := make([]Color, len(b.points))
	copy(newPoints, b.points)

	// check if all groups are proper
	newGroups := make(map[Point]*Group, len(b.groups))
	groupCopies := make(map[int]*Group) // map from group ID to copied group

	for k, v := range b.groups {
		if v != nil {
			// check if group is already copied
			if copiedGroup, exists := groupCopies[v.ID]; exists {
				newGroups[k] = copiedGroup
			} else {
				// no copy exists, make copy
				copiedGroup := v.copy()
				groupCopies[v.ID] = copiedGroup
				newGroups[k] = copiedGroup
			}
		}
	}

	newHistory := make([]uint64, len(b.history))
	copy(newHistory, b.history)

	return &Board{
		points:       newPoints,
		size:         b.size,
		internalSize: b.internalSize,
		groups:       newGroups,
		dsu:          b.dsu.copy(),
		history:      newHistory,
		koPoint:      b.koPoint,
		koHash:       b.koHash,
		nextGroupID:  b.nextGroupID,
	}
}

// convert 1-based (x, y) coords to a Point
func (b *Board) ToPoint(x int, y int) Point {
	return Point(y*b.internalSize + x)
}

// returns board size (9, 13, 19, etc)
func (b *Board) Size() int {
	return b.size
}

// returns the color at the given coords (1-based)
func (b *Board) At(x, y int) Color {
	p := b.ToPoint(x, y)
	return b.points[p]
}

// convert a Point to 1-based (x, y) coords
func (b *Board) ToXY(p Point) (int, int) {
	if b.internalSize == 0 {
		return 0, 0
	}
	y := int(p) / b.internalSize
	x := int(p) % b.internalSize
	return x, y
}

// returns the 4 direct neighbors of a point
func (b *Board) Neighbors(p Point) [4]Point {
	internalSize := Point(b.internalSize)
	return [4]Point{
		p - 1,            // left
		p + 1,            // right
		p - internalSize, // top
		p + internalSize, // bottom
	}
}

// returns a string representation of the board for display
func (b *Board) String() string {
	var sb strings.Builder
	for y := 1; y <= b.size; y++ {
		for x := 1; x <= b.size; x++ {
			p := b.ToPoint(x, y)
			switch b.points[p] {
			case Empty:
				sb.WriteString(". ")
			case Black:
				sb.WriteString("X ")
			case White:
				sb.WriteString("O ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// validates and applies a move, returning a NEW board state
func (b *Board) ApplyMove(move Move) (*Board, error) {
	// validate move
	if err := b.validatePlacement(move); err != nil {
		return nil, err
	}

	// create a new board state by copying the current one
	newBoard := b.copy()

	// place stone
	newBoard.points[move.Point] = move.Color

	// create a new group for the placed stone
	newGroup := newBoard.createNewGroup(move.Point, move.Color)

	// merge with friendly neighbor groups
	newBoard.mergeFriendlyNeighbors(move.Point, newGroup)

	// update enemy neighbor liberties
	newBoard.updateEnemyLiberties(move.Point, move.Color)

	// resolve captures for any enemy groups now at 0 liberties
	capturedStones := newBoard.resolveCaptures(move.Point, move.Color)
	_ = capturedStones // probably will use later for Ko logic

	// check for suicide (after captures are resolved)
	if err := newBoard.validateSuicide(move.Point); err != nil {
		return nil, err
	}

	// update board hash and check for Ko
	newBoard.koHash = newBoard.computeHash()
	if newBoard.isPositionRepeated() {
		return nil, errors.New("illegal move: Ko rule violation")
	}

	// add current hash to history
	newBoard.history = append(newBoard.history, newBoard.koHash)

	return newBoard, nil
}
