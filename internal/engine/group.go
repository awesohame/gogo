package engine

// Group is a chain of connected stones of the same color
type Group struct {
	ID        int
	Stones    map[Point]struct{}
	Liberties map[Point]struct{}
	Color     Color
}

// creates new group for a single stone
func newGroup(id int, point Point, color Color, liberties map[Point]struct{}) *Group {
	stones := make(map[Point]struct{})
	stones[point] = struct{}{}

	return &Group{
		ID:        id,
		Stones:    stones,
		Liberties: liberties,
		Color:     color,
	}
}

// copy creates a deep copy of the group
func (g *Group) copy() *Group {
	newStones := make(map[Point]struct{}, len(g.Stones))
	for s := range g.Stones {
		newStones[s] = struct{}{}
	}

	newLiberties := make(map[Point]struct{}, len(g.Liberties))
	for l := range g.Liberties {
		newLiberties[l] = struct{}{}
	}

	return &Group{
		ID:        g.ID,
		Stones:    newStones,
		Liberties: newLiberties,
		Color:     g.Color,
	}
}

// combines the stones and liberties of another group into this one. other group should be discarded after merge
func (g *Group) mergeWith(other *Group) {
	// add all stones from other group
	for stone := range other.Stones {
		g.Stones[stone] = struct{}{}
	}

	// add liberties from other group
	for liberty := range other.Liberties {
		g.Liberties[liberty] = struct{}{}
	}

	// remove any liberties that are now occupied by stones in the merged group
	for stone := range g.Stones {
		delete(g.Liberties, stone)
	}
}

// inits a new group for a stone, calc its initial liberties, and adds it to the board
func (b *Board) createNewGroup(p Point, c Color) *Group {
	// calc initial liberties for the new stone
	liberties := b.calculateInitialLiberties(p)

	// create a new group for this stone
	b.nextGroupID++
	newGroup := newGroup(b.nextGroupID, p, c, liberties)
	b.groups[p] = newGroup
	return newGroup
}

// finds and merges the new group with any adjacent friendly groups
func (b *Board) mergeFriendlyNeighbors(p Point, newGroup *Group) {
	for _, n := range b.Neighbors(p) {
		if b.points[n] == newGroup.Color {
			neighborGroup := b.groups[n]
			if neighborGroup != nil && b.dsu.Find(newGroup.ID) != b.dsu.Find(neighborGroup.ID) {
				// find actual root groups by looking them up through any stone in the group
				// use newly placed stone for newGroup and the neighbor stone for neighborGroup
				rootGroup := b.groups[p]
				neighborRootGroup := b.groups[n]

				// DSU union
				b.dsu.Union(newGroup.ID, neighborGroup.ID)
				newRootID := b.dsu.Find(newGroup.ID)

				var mergedGroup *Group
				// group with new root ID absorbs the other
				if newRootID == rootGroup.ID {
					rootGroup.mergeWith(neighborRootGroup)
					mergedGroup = rootGroup
				} else {
					neighborRootGroup.mergeWith(rootGroup)
					mergedGroup = neighborRootGroup
				}

				// update the groups map for all stones in the merged group to point to the new single root group
				for stonePoint := range mergedGroup.Stones {
					b.groups[stonePoint] = mergedGroup
				}
			}
		}
	}
}
