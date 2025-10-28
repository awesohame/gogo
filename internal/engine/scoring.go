package engine

// Score is final score for both players
type Score struct {
	Black       int
	White       int
	BlackStones int
	WhiteStones int
	BlackArea   int
	WhiteArea   int
	DamePoints  int
}

// computes the Chinese area score (stones + territory) for curr board state
// uses half-counting (total_points = black_score + white_score + dame)
func (b *Board) CalculateChineseScore() Score {
	score := Score{}
	visited := make(map[Point]bool)

	// count stones and territories
	for y := 1; y <= b.size; y++ {
		for x := 1; x <= b.size; x++ {
			p := b.ToPoint(x, y)

			if visited[p] {
				continue
			}

			color := b.points[p]

			switch color {
			case Black:
				score.BlackStones++
				visited[p] = true

			case White:
				score.WhiteStones++
				visited[p] = true

			case Empty:
				// find territory using flood fill
				territory, owner := b.floodFillTerritory(p, visited)
				territorySize := len(territory)

				switch owner {
				case Black:
					score.BlackArea += territorySize
				case White:
					score.WhiteArea += territorySize
				default:
					score.DamePoints += territorySize
				}
			}
		}
	}

	// final scores = stones + territory
	score.Black = score.BlackStones + score.BlackArea
	score.White = score.WhiteStones + score.WhiteArea

	return score
}

// computes score with komi (handicap)
func (b *Board) CalculateScoreWithKomi(komi float64) (black float64, white float64, winner Color) {
	score := b.CalculateChineseScore()

	blackScore := float64(score.Black)
	whiteScore := float64(score.White) + komi

	if blackScore > whiteScore {
		return blackScore, whiteScore, Black
	} else if whiteScore > blackScore {
		return blackScore, whiteScore, White
	}
	return blackScore, whiteScore, Empty // draw
}

// uses BFS to find connected empty territory and determine its owner
// returns territory points and color that owns it (or Empty if neutral/dame)
func (b *Board) floodFillTerritory(start Point, visited map[Point]bool) ([]Point, Color) {
	territory := []Point{}
	queue := []Point{start}
	visited[start] = true

	var owner Color = Empty
	bordersBlack := false
	bordersWhite := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		territory = append(territory, current)

		// check all neighbors
		for _, n := range b.Neighbors(current) {
			neighborColor := b.points[n]

			if neighborColor == Empty && !visited[n] {
				// continue flood fill into empty neighbor
				visited[n] = true
				queue = append(queue, n)
			} else if neighborColor == Black {
				bordersBlack = true
			} else if neighborColor == White {
				bordersWhite = true
			}
		}
	}

	// determine owner based on which colors border this territory
	if bordersBlack && !bordersWhite {
		owner = Black
	} else if bordersWhite && !bordersBlack {
		owner = White
	} else {
		// borders both colors or neither - neutral territory (dame)
		owner = Empty
	}

	return territory, owner
}
