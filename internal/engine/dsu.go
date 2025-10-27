package engine

// DSU
type DSU struct {
	parent []int
	rank   []int
}

// NewDSU creates a new DSU structure of a given size
func NewDSU(size int) *DSU {
	parent := make([]int, size)
	rank := make([]int, size) // rank is initialized to 0
	for i := 0; i < size; i++ {
		parent[i] = i
	}
	return &DSU{parent: parent, rank: rank}
}

// copy creates a deep copy of the DSU structure
func (d *DSU) copy() *DSU {
	newParent := make([]int, len(d.parent))
	newRank := make([]int, len(d.rank))
	copy(newParent, d.parent)
	copy(newRank, d.rank)
	return &DSU{parent: newParent, rank: newRank}
}

// find with path compression
func (d *DSU) Find(i int) int {
	if d.parent[i] == i {
		return i
	}
	// path compression
	d.parent[i] = d.Find(d.parent[i])
	return d.parent[i]
}

// Union by rank
func (d *DSU) Union(i, j int) {
	rootI := d.Find(i)
	rootJ := d.Find(j)
	if rootI != rootJ {
		if d.rank[rootI] < d.rank[rootJ] {
			d.parent[rootI] = rootJ
		} else if d.rank[rootJ] < d.rank[rootI] {
			d.parent[rootJ] = rootI
		} else {
			d.parent[rootJ] = rootI
			d.rank[rootI]++
		}
	}
}
