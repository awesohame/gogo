package engine

// returns set of points that contain dead stones
func (b *Board) findDeadStones() map[Point]bool {
	deadStones := make(map[Point]bool)
	analyzed := make(map[Point]bool)

	// analyze each group once
	for y := 1; y <= b.size; y++ {
		for x := 1; x <= b.size; x++ {
			p := b.ToPoint(x, y)
			color := b.points[p]

			if color == Empty || color == Border || analyzed[p] {
				continue
			}

			// get the group this stone belongs to
			group, ok := b.groups[p]
			if !ok {
				continue
			}

			// skip if already analyzed this group
			alreadyChecked := false
			for stone := range group.Stones {
				if analyzed[stone] {
					alreadyChecked = true
					break
				}
			}
			if alreadyChecked {
				continue
			}

			// mark all stones in group as analyzed
			for stone := range group.Stones {
				analyzed[stone] = true
			}

			// check if group is dead
			if b.isGroupDead(group) {
				for stone := range group.Stones {
					deadStones[stone] = true
				}
			}
		}
	}

	return deadStones
}

// checks if a group is dead (<2 eyes or surrounded)
func (b *Board) isGroupDead(group *Group) bool {
	// if group has 2+ liberties, likely alive
	if len(group.Liberties) >= 2 {
		return false
	}

	// if in atari (1 liberty), check if it's in enemy territory
	if len(group.Liberties) == 1 {
		return b.isInEnemyTerritory(group)
	}

	// no liberties = already captured
	return true
}

// checks if a group is completely surrounded by enemy stones
func (b *Board) isInEnemyTerritory(group *Group) bool {
	enemyColor := White
	if group.Color == White {
		enemyColor = Black
	}

	// if any liberty connects to friendly territory or has space, not dead
	for liberty := range group.Liberties {
		// check neighbors of the liberty
		friendlyInfluence := 0
		enemyInfluence := 0
		emptySpaces := 0

		for _, n := range b.Neighbors(liberty) {
			nColor := b.points[n]
			switch nColor {
			case group.Color:
				friendlyInfluence++
			case enemyColor:
				enemyInfluence++
			case Empty:
				emptySpaces++
			}
		}

		// if liberty has space to expand or friendly support, not dead
		if emptySpaces > 0 || friendlyInfluence > enemyInfluence {
			return false
		}
	}

	// completely surrounded by enemy
	return true
}

// removes dead stones from board for scoring
// returns a new board with dead stones removed
func (b *Board) removeDeadStones() *Board {
	deadStones := b.findDeadStones()

	if len(deadStones) == 0 {
		return b
	}

	// create new board without dead stones
	newBoard := b.copy()

	for deadPoint := range deadStones {
		newBoard.points[deadPoint] = Empty
		delete(newBoard.groups, deadPoint)
	}

	// rebuild groups for remaining stones
	newBoard.rebuildGroups()

	return newBoard
}

// rebuilds group structure after removing dead stones
func (b *Board) rebuildGroups() {
	b.groups = make(map[Point]*Group)
	b.dsu = NewDSU(b.size * b.size * 2)
	b.nextGroupID = 0

	// place all stones again to rebuild groups
	for y := 1; y <= b.size; y++ {
		for x := 1; x <= b.size; x++ {
			p := b.ToPoint(x, y)
			color := b.points[p]

			if color == Black || color == White {
				// create group for this stone
				liberties := b.calculateInitialLiberties(p)
				group := newGroup(b.nextGroupID, p, color, liberties)
				b.nextGroupID++
				b.groups[p] = group

				// merge with adjacent friendly groups
				b.mergeFriendlyNeighbors(p, group)
			}
		}
	}
}
