package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func NewGameState(name string, maxPlayers int, seed int64) *GameState {
	return &GameState{
		ID:             uuid.New().String(),
		Name:           name,
		MaxPlayers:     maxPlayers,
		CurrentTurn:    0,
		MaxTurns:       50,
		Phase:          PhasePlanning,
		Status:         "waiting",
		Map:            GenerateMap(seed, maxPlayers),
		Players:        make(map[string]*PlayerState),
		TurnOrder:      make([]string, 0),
		Events:         make([]GameEvent, 0),
		PendingActions: make([]PlayerAction, 0),
		LastUpdate:     time.Now(),
	}
}

func (gs *GameState) AddPlayer(playerID, username string, race Race, color string, isHost bool) {
	player := &PlayerState{
		PlayerID:       playerID,
		Username:       username,
		Race:           race,
		Color:          color,
		IsHost:         isHost,
		Resources:      ResourceStorage{Resources: Resources{Stone: 100, Metal: 50, GlowMushroom: 80, Water: 60, MagicCrystal: 10, FossilFuel: 20}, Capacity: 500},
		Production:     ResourceProduction{},
		Population:     10,
		PopulationCap:  20,
		ResearchPoints: 0,
		Techs:          make(map[string]bool),
		ResearchQueue:  make([]string, 0),
		Units:          make([]Unit, 0),
		Buildings:      make(map[string]*Building),
		VisionRange:    3,
		Diplomacy:      make(map[string]*DiplomaticRelation),
		TradeOffers:    make([]TradeOffer, 0),
		Embargoes:      make([]Embargo, 0),
		Ready:          false,
	}

	raceInfo := RaceDefs[race]
	if bonus, ok := raceInfo.Bonuses["vision_range"]; ok {
		player.VisionRange += int(bonus)
	}

	gs.Players[playerID] = player
	gs.TurnOrder = append(gs.TurnOrder, playerID)
}

func (gs *GameState) StartGame() {
	if len(gs.Players) < 1 {
		return
	}

	gs.Status = "playing"
	gs.Phase = PhasePlanning
	gs.CurrentTurn = 1

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	startPositions := FindPlayerStartPositions(gs.Map, len(gs.Players), r)

	idx := 0
	for _, playerID := range gs.TurnOrder {
		player := gs.Players[playerID]
		if idx >= len(startPositions) {
			idx = len(startPositions) - 1
		}
		pos := startPositions[idx]
		idx++

		layer := gs.Map.Layers[pos.Layer]
		cell, _ := GetCell(layer, pos.Coord)
		cell.Owner = playerID
		cell.Discovered = true

		building := &Building{
			Type:      BuildingMainBase,
			HP:        500,
			MaxHP:     500,
			Level:     1,
			Owner:     playerID,
			Completed: true,
			BuildTime: 0,
		}
		cell.Building = building

		buildingKey := fmt.Sprintf("%d_%s", pos.Layer, HexCoordKey(pos.Coord))
		player.Buildings[buildingKey] = building

		for i := 0; i < 3; i++ {
			unit := Unit{
				ID:      uuid.New().String(),
				Type:    UnitSapper,
				Owner:   playerID,
				HP:      30,
				MaxHP:   30,
				Attack:  5,
				Defense: 3,
				Range:   1,
				Speed:   2,
				Layer:   pos.Layer,
				Coord:   pos.Coord,
			}
			player.Units = append(player.Units, unit)
			cell.Units = append(cell.Units, unit)
		}

		gs.revealArea(player, pos.Layer, pos.Coord, player.VisionRange)
		gs.calculateProduction(player)
	}

	gs.PlanningEndsAt = time.Now().Add(90 * time.Second)
}

func (gs *GameState) revealArea(player *PlayerState, layer int, center HexCoord, radius int) {
	l := gs.Map.Layers[layer]
	for q := -radius; q <= radius; q++ {
		for r := -radius; r <= radius; r++ {
			coord := HexCoord{Q: center.Q + q, R: center.R + r}
			if HexDistance(center, coord) <= radius {
				if cell, ok := GetCell(l, coord); ok {
					cell.Discovered = true
				}
			}
		}
	}
}

func (gs *GameState) calculateProduction(player *PlayerState) {
	prod := ResourceProduction{}

	for _, building := range player.Buildings {
		if !building.Completed {
			continue
		}
		switch building.Type {
		case BuildingWorkshop:
			prod.Metal += 2
		case BuildingFungusFarm:
			prod.GlowMushroom += 3
		case BuildingAcademy:
			prod.MagicCrystal += 0
		}
	}

	for _, layer := range gs.Map.Layers {
		for _, cell := range layer.Cells {
			if cell.Owner == player.PlayerID && cell.MineralType != "none" {
				switch cell.MineralType {
				case "stone":
					prod.Stone += 2
				case "iron", "copper":
					prod.Metal += 2
				case "gold":
					prod.Metal += 3
				case "glow_mushroom":
					prod.GlowMushroom += 2
				case "magic_crystal":
					prod.MagicCrystal += 1
				case "fossil_fuel":
					prod.FossilFuel += 2
				}
			}
		}
	}

	raceInfo := RaceDefs[player.Race]
	if bonus, ok := raceInfo.Bonuses["metal_production"]; ok {
		prod.Metal = int(float64(prod.Metal) * (1 + bonus))
	}
	if bonus, ok := raceInfo.Bonuses["fungus_production"]; ok {
		prod.GlowMushroom = int(float64(prod.GlowMushroom) * (1 + bonus))
	}

	for techID := range player.Techs {
		tech := getTechByID(techID)
		if tech != nil {
			if bonus, ok := tech.Effects["metal_production"]; ok {
				prod.Metal = int(float64(prod.Metal) * (1 + bonus))
			}
			if bonus, ok := tech.Effects["fungus_production"]; ok {
				prod.GlowMushroom = int(float64(prod.GlowMushroom) * (1 + bonus))
			}
			if bonus, ok := tech.Effects["all_production"]; ok {
				prod.Stone = int(float64(prod.Stone) * (1 + bonus))
				prod.Metal = int(float64(prod.Metal) * (1 + bonus))
				prod.GlowMushroom = int(float64(prod.GlowMushroom) * (1 + bonus))
				prod.Water = int(float64(prod.Water) * (1 + bonus))
				prod.MagicCrystal = int(float64(prod.MagicCrystal) * (1 + bonus))
				prod.FossilFuel = int(float64(prod.FossilFuel) * (1 + bonus))
			}
		}
	}

	player.Production = prod
}

func getTechByID(id string) *Tech {
	for i := range TechDefs {
		if TechDefs[i].ID == id {
			return &TechDefs[i]
		}
	}
	return nil
}

func (gs *GameState) ProcessTurn() {
	if gs.Phase != PhasePlanning {
		return
	}

	gs.Phase = PhaseExecuting

	gs.processMining()
	gs.processBuilding()
	gs.processProduction()
	gs.processResearch()
	gs.processUnitMovement()
	gs.processCombat()
	gs.processDiplomacy()
	gs.processEvents()

	gs.CurrentTurn++
	gs.checkVictoryConditions()

	if gs.Phase != PhaseEnded {
		gs.Phase = PhasePlanning
		gs.PlanningEndsAt = time.Now().Add(90 * time.Second)
		gs.resetUnitActions()
		gs.PendingActions = make([]PlayerAction, 0)
	}

	gs.LastUpdate = time.Now()
}

func (gs *GameState) processProduction() {
	for _, player := range gs.Players {
		if player.Eliminated {
			continue
		}
		gs.calculateProduction(player)

		player.Resources.Stone = min(player.Resources.Stone+player.Production.Stone, player.Resources.Capacity)
		player.Resources.Metal = min(player.Resources.Metal+player.Production.Metal, player.Resources.Capacity)
		player.Resources.GlowMushroom = min(player.Resources.GlowMushroom+player.Production.GlowMushroom, player.Resources.Capacity)
		player.Resources.Water = min(player.Resources.Water+player.Production.Water, player.Resources.Capacity)
		player.Resources.MagicCrystal = min(player.Resources.MagicCrystal+player.Production.MagicCrystal, player.Resources.Capacity)
		player.Resources.FossilFuel = min(player.Resources.FossilFuel+player.Production.FossilFuel, player.Resources.Capacity)

		researchPts := gs.calculateResearchPoints(player)
		player.ResearchPoints += researchPts
	}
}

func (gs *GameState) calculateResearchPoints(player *PlayerState) int {
	points := 0

	for _, building := range player.Buildings {
		if building.Completed && building.Type == BuildingAcademy {
			points++
		}
	}

	for techID := range player.Techs {
		tech := getTechByID(techID)
		if tech != nil {
			if bonus, ok := tech.Effects["research_point"]; ok {
				points += int(bonus)
			}
		}
	}

	return points
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (gs *GameState) processMining() {
	for _, action := range gs.PendingActions {
		if action.Action != "mine" {
			continue
		}

		player := gs.Players[action.PlayerID]
		if player == nil || player.Eliminated {
			continue
		}

		layerIdx := int(action.Data["layer"].(float64))
		coordData := action.Data["coord"].(map[string]interface{})
		coord := HexCoord{Q: int(coordData["q"].(float64)), R: int(coordData["r"].(float64))}

		if layerIdx < 0 || layerIdx >= len(gs.Map.Layers) {
			continue
		}

		layer := gs.Map.Layers[layerIdx]
		cell, ok := GetCell(layer, coord)
		if !ok || !cell.IsWall {
			continue
		}

		if cell.MiningOwner != "" {
			continue
		}

		canMine := false
		for _, n := range getHexNeighbors(coord) {
			if nc, ok := GetCell(layer, n); ok && nc.Owner == action.PlayerID {
				canMine = true
				break
			}
		}
		if !canMine {
			continue
		}

		cell.MiningOwner = action.PlayerID
		cell.MiningProgress = 0
		cell.MiningTotal = cell.RockHardness * 2
	}

	for layerIdx, layer := range gs.Map.Layers {
		for _, cell := range layer.Cells {
			if cell.MiningOwner == "" || !cell.IsWall {
				continue
			}

			player := gs.Players[cell.MiningOwner]
			if player == nil || player.Eliminated {
				cell.MiningOwner = ""
				cell.MiningProgress = 0
				cell.MiningTotal = 0
				continue
			}

			miningSpeed := 1.0
			raceInfo := RaceDefs[player.Race]
			if bonus, ok := raceInfo.Bonuses["mining_speed"]; ok {
				miningSpeed *= bonus
			}

			if player.Techs["basic_mining"] {
				miningSpeed *= 1.2
			}
			if player.Techs["advanced_mining"] {
				miningSpeed *= 1.3
			}

			cell.MiningProgress += int(miningSpeed)
			if cell.MiningProgress >= cell.MiningTotal {
				collapseChance := cell.WaterContent * 0.1
				if player.Techs["support_beams"] {
					collapseChance *= 0.5
				}

				if rand.Float64() < collapseChance {
					cell.MiningOwner = ""
					cell.MiningProgress = 0
					cell.MiningTotal = 0
					gs.addEvent(GameEvent{
						Type:     EventEarthquake,
						Message:  "塌方！挖掘失败，损失了一些资源",
						Turn:     gs.CurrentTurn,
						Layer:    layerIdx,
						Location: cell.Coord,
					})
					player.Resources.Stone = max(0, player.Resources.Stone-10)
					continue
				}

				cell.IsWall = false
				cell.Owner = cell.MiningOwner
				cell.Discovered = true
				cell.MiningOwner = ""
				cell.MiningProgress = 0
				cell.MiningTotal = 0

				gs.revealArea(player, layerIdx, cell.Coord, 2)
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (gs *GameState) processBuilding() {
	for _, action := range gs.PendingActions {
		if action.Action != "build" {
			continue
		}

		player := gs.Players[action.PlayerID]
		if player == nil || player.Eliminated {
			continue
		}

		layerIdx := int(action.Data["layer"].(float64))
		coordData := action.Data["coord"].(map[string]interface{})
		coord := HexCoord{Q: int(coordData["q"].(float64)), R: int(coordData["r"].(float64))}
		buildingType := BuildingType(action.Data["building"].(string))

		if layerIdx < 0 || layerIdx >= len(gs.Map.Layers) {
			continue
		}

		layer := gs.Map.Layers[layerIdx]
		cell, ok := GetCell(layer, coord)
		if !ok || cell.IsWall || cell.Building != nil || cell.Owner != action.PlayerID {
			continue
		}

		buildingDef, ok := BuildingDefs[buildingType]
		if !ok {
			continue
		}

		if buildingDef.RaceSpecific != "" && Race(buildingDef.RaceSpecific) != player.Race {
			continue
		}

		if !gs.canAfford(player, buildingDef.Cost) {
			continue
		}

		gs.payCost(player, buildingDef.Cost)

		building := &Building{
			Type:      buildingType,
			HP:        buildingDef.HP,
			MaxHP:     buildingDef.HP,
			Level:     1,
			Owner:     action.PlayerID,
			Completed: false,
			Progress:  0,
			BuildTime: buildingDef.BuildTime,
		}

		hpBonus := 1.0
		if bonus, ok := RaceDefs[player.Race].Bonuses["building_hp"]; ok {
			hpBonus = bonus
		}
		building.HP = int(float64(building.HP) * hpBonus)
		building.MaxHP = building.HP

		cell.Building = building
		buildingKey := fmt.Sprintf("%d_%s", layerIdx, HexCoordKey(coord))
		player.Buildings[buildingKey] = building
	}

	for _, player := range gs.Players {
		if player.Eliminated {
			continue
		}

		for _, building := range player.Buildings {
			if !building.Completed {
				buildSpeed := 1.0
				if player.Techs["stonecutting"] {
					buildSpeed += 0.25
				}

				building.Progress += int(float64(1) * buildSpeed)
				if building.Progress >= building.BuildTime {
					building.Completed = true
					building.Progress = building.BuildTime

					if building.Type == BuildingLivingQuarters {
						player.PopulationCap += 5
					}
					if building.Type == BuildingWarehouse {
						player.Resources.Capacity += 100
					}
				}
			}
		}
	}
}

func (gs *GameState) canAfford(player *PlayerState, cost Resources) bool {
	return player.Resources.Stone >= cost.Stone &&
		player.Resources.Metal >= cost.Metal &&
		player.Resources.GlowMushroom >= cost.GlowMushroom &&
		player.Resources.Water >= cost.Water &&
		player.Resources.MagicCrystal >= cost.MagicCrystal &&
		player.Resources.FossilFuel >= cost.FossilFuel
}

func (gs *GameState) payCost(player *PlayerState, cost Resources) {
	player.Resources.Stone -= cost.Stone
	player.Resources.Metal -= cost.Metal
	player.Resources.GlowMushroom -= cost.GlowMushroom
	player.Resources.Water -= cost.Water
	player.Resources.MagicCrystal -= cost.MagicCrystal
	player.Resources.FossilFuel -= cost.FossilFuel
}

func (gs *GameState) processResearch() {
	for _, player := range gs.Players {
		if player.Eliminated {
			continue
		}

		if player.CurrentResearch == "" {
			if len(player.ResearchQueue) > 0 {
				player.CurrentResearch = player.ResearchQueue[0]
				player.ResearchQueue = player.ResearchQueue[1:]
			}
			continue
		}

		tech := getTechByID(player.CurrentResearch)
		if tech == nil {
			player.CurrentResearch = ""
			continue
		}

		researchSpeed := 1.0
		if player.Techs["scientific_method"] {
			researchSpeed += 0.3
		}
		if player.Techs["crystal_resonance"] {
			researchSpeed += 0.2
		}

		player.ResearchProgress += int(float64(player.ResearchPoints) * researchSpeed)

		if player.ResearchProgress >= tech.Cost {
			player.Techs[player.CurrentResearch] = true
			player.ResearchProgress = 0
			player.CurrentResearch = ""

			if len(player.ResearchQueue) > 0 {
				player.CurrentResearch = player.ResearchQueue[0]
				player.ResearchQueue = player.ResearchQueue[1:]
			}
		}
	}
}

func (gs *GameState) processUnitMovement() {
	for _, action := range gs.PendingActions {
		if action.Action != "move" {
			continue
		}

		player := gs.Players[action.PlayerID]
		if player == nil || player.Eliminated {
			continue
		}

		unitID := action.Data["unit_id"].(string)
		toData := action.Data["to"].(map[string]interface{})
		to := HexCoord{Q: int(toData["q"].(float64)), R: int(toData["r"].(float64))}
		toLayer := int(action.Data["to_layer"].(float64))

		unit := gs.findUnit(unitID)
		if unit == nil || unit.Owner != action.PlayerID || unit.Moved {
			continue
		}

		if toLayer != unit.Layer {
			layer := gs.Map.Layers[unit.Layer]
			cell, _ := GetCell(layer, unit.Coord)
			if !cell.IsShaft && cell.Building == nil {
				continue
			}
			if cell.Building != nil && cell.Building.Type != BuildingElevator {
				continue
			}
		}

		distance := HexDistance(unit.Coord, to)
		if distance > unit.Speed {
			continue
		}

		if toLayer < 0 || toLayer >= len(gs.Map.Layers) {
			continue
		}
		layer := gs.Map.Layers[toLayer]
		cell, ok := GetCell(layer, to)
		if !ok || cell.IsWall {
			continue
		}

		gs.removeUnitFromCell(unit)
		unit.Coord = to
		unit.Layer = toLayer
		unit.Moved = true
		cell.Units = append(cell.Units, *unit)

		gs.revealArea(player, toLayer, to, 2)
	}
}

func (gs *GameState) findUnit(unitID string) *Unit {
	for _, player := range gs.Players {
		for i := range player.Units {
			if player.Units[i].ID == unitID {
				return &player.Units[i]
			}
		}
	}
	return nil
}

func (gs *GameState) removeUnitFromCell(unit *Unit) {
	if unit.Layer < 0 || unit.Layer >= len(gs.Map.Layers) {
		return
	}
	layer := gs.Map.Layers[unit.Layer]
	cell, ok := GetCell(layer, unit.Coord)
	if !ok {
		return
	}

	newUnits := make([]Unit, 0, len(cell.Units)-1)
	for _, u := range cell.Units {
		if u.ID != unit.ID {
			newUnits = append(newUnits, u)
		}
	}
	cell.Units = newUnits
}

func (gs *GameState) resetUnitActions() {
	for _, player := range gs.Players {
		for i := range player.Units {
			player.Units[i].Moved = false
			player.Units[i].Attacked = false
		}
	}

	for _, layer := range gs.Map.Layers {
		for _, cell := range layer.Cells {
			for i := range cell.Units {
				cell.Units[i].Moved = false
				cell.Units[i].Attacked = false
			}
		}
	}
}

func (gs *GameState) addEvent(event GameEvent) {
	gs.Events = append(gs.Events, event)
	if len(gs.Events) > 50 {
		gs.Events = gs.Events[1:]
	}
}

func (gs *GameState) processEvents() {
	if gs.CurrentTurn%5 != 0 {
		return
	}

	eventTypes := []GameEventType{
		EventEarthquake,
		EventFlood,
		EventMineralVein,
		EventAncientRuin,
	}

	eventType := eventTypes[rand.Intn(len(eventTypes))]
	layer := rand.Intn(5)

	switch eventType {
	case EventEarthquake:
		gs.addEvent(GameEvent{
			Type:     EventEarthquake,
			Message:  "地震！某区域建筑损失10%耐久",
			Turn:     gs.CurrentTurn,
			Layer:    layer,
			Duration: 0,
		})
		l := gs.Map.Layers[layer]
		damageCount := 0
		for _, cell := range l.Cells {
			if cell.Building != nil && damageCount < 10 {
				cell.Building.HP = int(float64(cell.Building.HP) * 0.9)
				damageCount++
			}
		}

	case EventMineralVein:
		gs.addEvent(GameEvent{
			Type:     EventMineralVein,
			Message:  "发现新矿脉！",
			Turn:     gs.CurrentTurn,
			Layer:    layer,
			Duration: 0,
		})
		l := gs.Map.Layers[layer]
		found := 0
		for _, cell := range l.Cells {
			if !cell.IsWall && cell.MineralType == "none" && found < 5 {
				minerals := []string{"iron", "copper", "gold", "magic_crystal", "fossil_fuel"}
				cell.MineralType = minerals[rand.Intn(len(minerals))]
				found++
			}
		}
	}
}

func (gs *GameState) checkVictoryConditions() {
	alivePlayers := make([]string, 0)
	for _, playerID := range gs.TurnOrder {
		player := gs.Players[playerID]
		if !player.Eliminated {
			hasMainBase := false
			for _, building := range player.Buildings {
				if building.Type == BuildingMainBase && building.Completed && building.HP > 0 {
					hasMainBase = true
					break
				}
			}
			if !hasMainBase {
				player.Eliminated = true
			} else {
				alivePlayers = append(alivePlayers, playerID)
			}
		}
	}

	if len(alivePlayers) == 1 {
		gs.WinnerID = alivePlayers[0]
		gs.VictoryType = "conquest"
		gs.Phase = PhaseEnded
		gs.Status = "finished"
		return
	}

	for _, playerID := range gs.TurnOrder {
		player := gs.Players[playerID]
		if len(player.Techs) >= len(TechDefs) {
			gs.WinnerID = playerID
			gs.VictoryType = "tech"
			gs.Phase = PhaseEnded
			gs.Status = "finished"
			return
		}
	}

	if gs.CurrentTurn >= gs.MaxTurns {
		bestScore := 0
		bestPlayer := ""
		for _, playerID := range gs.TurnOrder {
			player := gs.Players[playerID]
			if player.Eliminated {
				continue
			}

			territoryCount := 0
			buildingCount := 0
			militaryPower := 0

			for _, layer := range gs.Map.Layers {
				for _, cell := range layer.Cells {
					if cell.Owner == playerID {
						territoryCount++
					}
				}
			}

			for _, building := range player.Buildings {
				if building.Completed {
					buildingCount++
				}
			}

			militaryPower = len(player.Units) * 10

			score := territoryCount*1 + buildingCount*2 + int(float64(militaryPower)*1.5) + len(player.Techs)*3
			player.Score = score

			if score > bestScore {
				bestScore = score
				bestPlayer = playerID
			}
		}

		gs.WinnerID = bestPlayer
		gs.VictoryType = "score"
		gs.Phase = PhaseEnded
		gs.Status = "finished"
	}
}
