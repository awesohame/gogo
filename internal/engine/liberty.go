package engine

// calc liberties for a newly placed stone
func (b *Board) calculateInitialLiberties(p Point) map[Point]struct{} {
	liberties := make(map[Point]struct{})
	for _, n := range b.Neighbors(p) {
		if b.points[n] == Empty {
			liberties[n] = struct{}{}
		}
	}
	return liberties
}

// updates the liberty counts of enemy groups adjacent to a newly placed stone
func (b *Board) updateEnemyLiberties(p Point, placedColor Color) {
	for _, n := range b.Neighbors(p) {
		if b.points[n] != Empty && b.points[n] != placedColor {
			// this is an enemy neighbor
			enemyGroup := b.groups[n]
			if enemyGroup != nil {
				// a liberty is taken away from this group
				delete(enemyGroup.Liberties, p)
			}
		}
	}
}
