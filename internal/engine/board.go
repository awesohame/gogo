package engine

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
		history:      make([]uint64, 0),
		koPoint:      -1, // use -1 for no active Ko point
	}
}
