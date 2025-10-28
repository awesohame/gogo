package engine

import "errors"

// Move is a player's action, placing a stone on the board
type Move struct {
	Point Point
	Color Color
}

// check if a move is on an empty point
func (b *Board) validatePlacement(m Move) error {
	if b.points[m.Point] != Empty {
		return errors.New("point is not empty")
	}
	return nil
}

// check if a move is suicidal (creates a group with zero liberties without capturing)
// this is called AFTER captures have been resolved
func (b *Board) validateSuicide(p Point) error {
	group := b.groups[p]
	if group != nil && len(group.Liberties) == 0 {
		return errors.New("suicidal move: group has no liberties")
	}
	return nil
}
