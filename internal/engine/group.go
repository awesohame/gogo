package engine

// Group is a chain of connected stones of the same color
type Group struct {
	ID        int
	Stones    map[Point]struct{}
	Liberties map[Point]struct{}
	Color     Color
}
