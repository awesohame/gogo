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
