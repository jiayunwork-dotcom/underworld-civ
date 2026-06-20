package game

import (
	"fmt"
	"math/rand"
)

const (
	LayerShallow    = 0
	LayerMiddle     = 1
	LayerDeep       = 2
	LayerLava       = 3
	LayerAncient    = 4
)

var LayerNames = []string{
	"浅层 - 蘑菇林带",
	"中层 - 石灰岩洞窟",
	"深层 - 金属矿脉",
	"熔岩层 - 魔晶之地",
	"远古层 - 化石遗迹",
}

func HexKey(q, r int) string {
	return fmt.Sprintf("%d,%d", q, r)
}

func HexCoordKey(c HexCoord) string {
	return HexKey(c.Q, c.R)
}

func GenerateMap(seed int64, playerCount int) *GameMap {
	r := rand.New(rand.NewSource(seed))

	gameMap := &GameMap{
		Layers: make([]*Layer, 5),
	}

	layerConfigs := []struct {
		width  int
		height int
		fillProb float64
	}{
		{20, 15, 0.45},
		{25, 20, 0.40},
		{30, 22, 0.42},
		{25, 18, 0.50},
		{20, 15, 0.48},
	}

	for i := 0; i < 5; i++ {
		cfg := layerConfigs[i]
		layer := generateLayer(r, i, cfg.width, cfg.height, cfg.fillProb)
		gameMap.Layers[i] = layer
	}

	placeShafts(gameMap, r)
	distributeResources(gameMap, r)

	return gameMap
}

func generateLayer(r *rand.Rand, depth, width, height int, fillProb float64) *Layer {
	layer := &Layer{
		Depth:  depth,
		Name:   LayerNames[depth],
		Width:  width,
		Height: height,
		Cells:  make(map[string]*HexCell),
	}

	for q := 0; q < width; q++ {
		for rIdx := 0; rIdx < height; rIdx++ {
			cell := &HexCell{
				Coord:        HexCoord{Q: q, R: rIdx},
				Layer:        depth,
				Discovered:   false,
				Owner:        "",
				RockHardness: 1 + r.Intn(5),
				WaterContent: r.Float64() * 0.8,
				MineralType:  "none",
				IsWall:       r.Float64() < fillProb,
				IsShaft:      false,
			}

			if q == 0 || q == width-1 || rIdx == 0 || rIdx == height-1 {
				cell.IsWall = true
				cell.RockHardness = 5
			}

			layer.Cells[HexKey(q, rIdx)] = cell
		}
	}

	for i := 0; i < 5; i++ {
		cellularAutomataStep(layer)
	}

	ensureConnectivity(layer, r)

	return layer
}

func cellularAutomataStep(layer *Layer) {
	newCells := make(map[string]*HexCell)

	for key, cell := range layer.Cells {
		neighbors := countWallNeighbors(layer, cell.Coord)

		newCell := *cell
		if cell.IsWall {
			if neighbors < 3 {
				newCell.IsWall = false
			}
		} else {
			if neighbors > 4 {
				newCell.IsWall = true
			}
		}
		newCells[key] = &newCell
	}

	layer.Cells = newCells
}

func countWallNeighbors(layer *Layer, coord HexCoord) int {
	count := 0
	neighbors := getHexNeighbors(coord)

	for _, n := range neighbors {
		key := HexKey(n.Q, n.R)
		if cell, ok := layer.Cells[key]; ok {
			if cell.IsWall {
				count++
			}
		} else {
			count++
		}
	}
	return count
}

func getHexNeighbors(coord HexCoord) []HexCoord {
	return []HexCoord{
		{Q: coord.Q + 1, R: coord.R},
		{Q: coord.Q - 1, R: coord.R},
		{Q: coord.Q, R: coord.R + 1},
		{Q: coord.Q, R: coord.R - 1},
		{Q: coord.Q + 1, R: coord.R - 1},
		{Q: coord.Q - 1, R: coord.R + 1},
	}
}

func ensureConnectivity(layer *Layer, r *rand.Rand) {
	visited := make(map[string]bool)
	var largestRegion map[string]bool

	for key, cell := range layer.Cells {
		if !cell.IsWall && !visited[key] {
			region := floodFill(layer, cell.Coord, visited)
			if len(region) > len(largestRegion) {
				largestRegion = region
			}
		}
	}

	for key, cell := range layer.Cells {
		if !cell.IsWall && !largestRegion[key] {
			cell.IsWall = true
		}
	}
}

func floodFill(layer *Layer, start HexCoord, visited map[string]bool) map[string]bool {
	region := make(map[string]bool)
	queue := []HexCoord{start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		key := HexKey(current.Q, current.R)
		if visited[key] {
			continue
		}

		cell, ok := layer.Cells[key]
		if !ok || cell.IsWall {
			continue
		}

		visited[key] = true
		region[key] = true

		for _, neighbor := range getHexNeighbors(current) {
			nKey := HexKey(neighbor.Q, neighbor.R)
			if !visited[nKey] {
				queue = append(queue, neighbor)
			}
		}
	}

	return region
}

func placeShafts(gameMap *GameMap, r *rand.Rand) {
	shaftCount := 3
	shaftLocations := make([]HexCoord, 0)

	layer0 := gameMap.Layers[0]
	for _, cell := range layer0.Cells {
		if !cell.IsWall && len(shaftLocations) < shaftCount {
			if r.Float64() < 0.1 {
				shaftLocations = append(shaftLocations, cell.Coord)
			}
		}
	}

	for len(shaftLocations) < shaftCount {
		for _, cell := range layer0.Cells {
			if !cell.IsWall {
				shaftLocations = append(shaftLocations, cell.Coord)
				break
			}
		}
	}

	for i := 0; i < 5; i++ {
		layer := gameMap.Layers[i]
		for _, loc := range shaftLocations {
			key := HexKey(loc.Q, loc.R)
			if cell, ok := layer.Cells[key]; ok {
				cell.IsWall = false
				cell.IsShaft = true
				cell.RockHardness = 1
			}
		}
		layer.Shafts = shaftLocations
	}
}

func distributeResources(gameMap *GameMap, r *rand.Rand) {
	mineralDist := []map[string]float64{
		{"glow_mushroom": 0.3, "stone": 0.2, "none": 0.5},
		{"stone": 0.3, "iron": 0.15, "copper": 0.1, "none": 0.45},
		{"iron": 0.25, "copper": 0.2, "gold": 0.05, "stone": 0.2, "none": 0.3},
		{"magic_crystal": 0.2, "gold": 0.1, "iron": 0.15, "none": 0.55},
		{"fossil_fuel": 0.25, "magic_crystal": 0.1, "gold": 0.05, "none": 0.6},
	}

	for i, layer := range gameMap.Layers {
		dist := mineralDist[i]
		for _, cell := range layer.Cells {
			if !cell.IsWall && !cell.IsShaft {
				roll := r.Float64()
				cumulative := 0.0
				for mineral, prob := range dist {
					cumulative += prob
					if roll < cumulative {
						cell.MineralType = mineral
						break
					}
				}
			}
		}
	}
}

func GetCell(layer *Layer, coord HexCoord) (*HexCell, bool) {
	key := HexKey(coord.Q, coord.R)
	cell, ok := layer.Cells[key]
	return cell, ok
}

func IsAdjacent(a, b HexCoord) bool {
	for _, n := range getHexNeighbors(a) {
		if n.Q == b.Q && n.R == b.R {
			return true
		}
	}
	return false
}

func HexDistance(a, b HexCoord) int {
	return (abs(a.Q-b.Q) + abs(a.Q+a.R-b.Q-b.R) + abs(a.R-b.R)) / 2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func FindPlayerStartPositions(gameMap *GameMap, playerCount int, r *rand.Rand) []struct {
	Layer int
	Coord HexCoord
} {
	positions := make([]struct {
		Layer int
		Coord HexCoord
	}, 0)

	layer := gameMap.Layers[0]
	candidates := make([]HexCoord, 0)

	for _, cell := range layer.Cells {
		if !cell.IsWall && !cell.IsShaft {
			neighbors := getHexNeighbors(cell.Coord)
			openCount := 0
			for _, n := range neighbors {
				if nc, ok := GetCell(layer, n); ok && !nc.IsWall {
					openCount++
				}
			}
			if openCount >= 4 {
				candidates = append(candidates, cell.Coord)
			}
		}
	}

	for i := 0; i < playerCount && len(candidates) > 0; i++ {
		idx := r.Intn(len(candidates))
		pos := candidates[idx]

		positions = append(positions, struct {
			Layer int
			Coord HexCoord
		}{Layer: 0, Coord: pos})

		newCandidates := make([]HexCoord, 0)
		for _, c := range candidates {
			if HexDistance(c, pos) > 5 {
				newCandidates = append(newCandidates, c)
			}
		}
		candidates = newCandidates
	}

	return positions
}
