package engine

import (
	"math/rand"
)

// zobrist hash table for position hashing
var zobristTable [19 * 19 * 2]uint64

func init() {
	// initialize Zobrist hash table with random values
	rng := rand.New(rand.NewSource(42)) // fixed seed
	for i := range zobristTable {
		zobristTable[i] = rng.Uint64()
	}
}

// compute Zobrist hash for curr board state
func (b *Board) computeHash() uint64 {
	hash := uint64(0)
	for y := 1; y <= b.size; y++ {
		for x := 1; x <= b.size; x++ {
			p := b.ToPoint(x, y)
			color := b.points[p]
			if color == Black || color == White {
				// compute idx = (point_idx * 2) + (0 for black, 1 for white)
				idx := (y-1)*b.size + (x - 1)
				if color == White {
					idx = idx*2 + 1
				} else {
					idx = idx * 2
				}
				hash ^= zobristTable[idx]
			}
		}
	}
	return hash
}

// check if the current position has occurred before (Ko or Superko)
func (b *Board) isPositionRepeated() bool {
	currentHash := b.koHash
	for _, pastHash := range b.history {
		if pastHash == currentHash {
			return true
		}
	}
	return false
}
