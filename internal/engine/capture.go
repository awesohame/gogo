package engine

// checks all enemy neighbors and removes groups with zero liberties
// returns count of captured stones
func (b *Board) resolveCaptures(p Point, placedColor Color) int {
	capturedCount := 0
	capturedGroups := make(map[int]bool) // track already processed groups

	for _, n := range b.Neighbors(p) {
		if b.points[n] != Empty && b.points[n] != placedColor {
			// this is an enemy neighbor
			enemyGroup := b.groups[n]
			if enemyGroup != nil && !capturedGroups[enemyGroup.ID] {
				// check if this group has zero liberties
				if len(enemyGroup.Liberties) == 0 {
					// capture this group
					capturedCount += b.captureGroup(enemyGroup)
					capturedGroups[enemyGroup.ID] = true
				}
			}
		}
	}

	return capturedCount
}

// removes all stones of a group from the board and updates liberties of adjacent groups
// returns no. of stones captured
func (b *Board) captureGroup(group *Group) int {
	capturedCount := 0

	// for each stone in the captured group
	for stone := range group.Stones {
		// remove the stone from the board
		b.points[stone] = Empty

		// remove from groups map
		delete(b.groups, stone)

		capturedCount++

		// update liberties of adjacent groups (they now have a new liberty)
		for _, n := range b.Neighbors(stone) {
			if b.points[n] != Empty {
				neighborGroup := b.groups[n]
				if neighborGroup != nil {
					// this adjacent group gains the newly empty point as a liberty
					neighborGroup.Liberties[stone] = struct{}{}
				}
			}
		}
	}

	return capturedCount
}
