package game

import (
	"fmt"
	"math"
	"math/rand"
)

func (gs *GameState) processCombat() {
	for _, action := range gs.PendingActions {
		if action.Action != "attack" {
			continue
		}

		player := gs.Players[action.PlayerID]
		if player == nil || player.Eliminated {
			continue
		}

		unitID := action.Data["unit_id"].(string)
		targetData := action.Data["target"].(map[string]interface{})
		target := HexCoord{Q: int(targetData["q"].(float64)), R: int(targetData["r"].(float64))}
		targetLayer := int(action.Data["target_layer"].(float64))

		unit := gs.findUnit(unitID)
		if unit == nil || unit.Owner != action.PlayerID || unit.Attacked {
			continue
		}

		distance := HexDistance(unit.Coord, target)
		if unit.Layer != targetLayer || distance > unit.Range {
			continue
		}

		layer := gs.Map.Layers[targetLayer]
		targetCell, ok := GetCell(layer, target)
		if !ok {
			continue
		}

		unit.Attacked = true

		attackPower := float64(unit.Attack)
		attackBonus := 1.0
		attackFlatBonus := 0.0

		raceInfo := RaceDefs[player.Race]
		racialBoost := player.GetTechEffect("racial_bonus_boost")

		if bonus, ok := raceInfo.Bonuses["attack"]; ok {
			attackBonus += bonus * (1 + racialBoost)
		}

		attackBonus += player.GetTechEffect("attack_pct")
		attackBonus += player.GetTechEffect("all_stats")
		attackFlatBonus += player.GetTechEffect("attack_flat")

		siegeBonus := player.GetTechEffect("siege_damage")
		if unit.Type == UnitSiegeRam {
			attackBonus += siegeBonus
		}

		attackPower = (attackPower + attackFlatBonus) * attackBonus

		hpPctBonus := player.GetTechEffect("hp_pct")
		unitMaxHP := float64(unit.MaxHP) * (1 + hpPctBonus + player.GetTechEffect("all_stats"))

		defenderPlayerID := ""
		if len(targetCell.Units) > 0 {
			defenderPlayerID = targetCell.Units[0].Owner
		} else if targetCell.Building != nil {
			defenderPlayerID = targetCell.Building.Owner
		}

		isAmbush := false
		if defenderPlayerID != "" && targetCell.Owner == defenderPlayerID {
			isAmbush = true
		}

		if len(targetCell.Units) > 0 {
			defenderUnit := &targetCell.Units[0]
			defender := gs.Players[defenderUnit.Owner]

			defensePower := float64(defenderUnit.Defense)
			defenseBonus := 1.0
			defenseFlatBonus := 0.0

			if isAmbush {
				defenseBonus += 0.5
			}

			if defender != nil {
				defenseBonus += defender.GetTechEffect("defense_pct")
				defenseBonus += defender.GetTechEffect("all_stats")
				defenseFlatBonus += defender.GetTechEffect("defense_flat")

				if defenderUnit.Type == UnitInfantry {
					defenseFlatBonus += defender.GetTechEffect("infantry_defense")
				}
			}

			defensePower = (defensePower + defenseFlatBonus) * defenseBonus

			damage := math.Max(1, attackPower-defensePower/2)
			damage = damage * (0.8 + rand.Float64()*0.4)

			_ = unitMaxHP

			defenderUnit.HP -= int(damage)

			gs.updateUnitInPlayerList(defenderUnit)

			if defenderUnit.HP <= 0 {
				gs.removeUnitFromCell(defenderUnit)
				gs.removeUnitFromPlayer(defenderUnit)
			}

			if distance <= 1 && !defenderUnit.Moved {
				counterAttack := float64(defenderUnit.Attack) * 0.5
				if defender != nil {
					counterAttackFlat := defender.GetTechEffect("attack_flat")
					counterAttackPct := 1 + defender.GetTechEffect("attack_pct") + defender.GetTechEffect("all_stats")
					counterAttack = (counterAttack + counterAttackFlat) * counterAttackPct
				}
				unitDefense := float64(unit.Defense)
				unitDefenseFlat := player.GetTechEffect("defense_flat")
				unitDefensePct := 1 + player.GetTechEffect("defense_pct") + player.GetTechEffect("all_stats")
				if unit.Type == UnitInfantry {
					unitDefenseFlat += player.GetTechEffect("infantry_defense")
				}
				unitDefense = (unitDefense + unitDefenseFlat) * unitDefensePct
				counterDamage := math.Max(1, counterAttack-unitDefense/2)

				unit.HP -= int(counterDamage)
				gs.updateUnitInPlayerList(unit)

				if unit.HP <= 0 {
					gs.removeUnitFromCell(unit)
					gs.removeUnitFromPlayer(unit)
				}
			}
		}

		if targetCell.Building != nil && targetCell.Building.Owner != action.PlayerID {
			buildingDamage := attackPower
			if unit.Type == UnitSiegeRam {
				buildingDamage *= 3
			}
			buildingDamage *= (1 + player.GetTechEffect("building_damage"))

			if targetCell.Building.Type == BuildingWall {
				buildingDamage *= 0.5
			}

			targetCell.Building.HP -= int(buildingDamage)

			if targetCell.Building.HP <= 0 {
				buildingOwner := targetCell.Building.Owner
				buildingKey := fmt.Sprintf("%d_%s", targetLayer, HexCoordKey(target))
				if ownerPlayer, ok := gs.Players[buildingOwner]; ok {
					delete(ownerPlayer.Buildings, buildingKey)
				}
				targetCell.Building = nil
			}
		}

		gs.checkMoraleAndRoute(targetLayer, target)
	}
}

func (gs *GameState) updateUnitInPlayerList(unit *Unit) {
	if player, ok := gs.Players[unit.Owner]; ok {
		for i := range player.Units {
			if player.Units[i].ID == unit.ID {
				player.Units[i] = *unit
				return
			}
		}
	}
}

func (gs *GameState) removeUnitFromPlayer(unit *Unit) {
	if player, ok := gs.Players[unit.Owner]; ok {
		newUnits := make([]Unit, 0, len(player.Units)-1)
		for _, u := range player.Units {
			if u.ID != unit.ID {
				newUnits = append(newUnits, u)
			}
		}
		player.Units = newUnits
	}
}

func (gs *GameState) checkMoraleAndRoute(layerIdx int, coord HexCoord) {
	layer := gs.Map.Layers[layerIdx]
	cell, ok := GetCell(layer, coord)
	if !ok {
		return
	}

	if len(cell.Units) == 0 {
		return
	}

	owner := cell.Units[0].Owner
	player := gs.Players[owner]
	if player == nil {
		return
	}

	totalUnits := 0
	lostUnits := 0

	nearbyUnits := make([]*Unit, 0)
	for _, n := range getHexNeighbors(coord) {
		if nc, ok := GetCell(layer, n); ok {
			for i := range nc.Units {
				if nc.Units[i].Owner == owner {
					nearbyUnits = append(nearbyUnits, &nc.Units[i])
				}
			}
		}
	}

	totalUnits = len(nearbyUnits)
	for _, u := range nearbyUnits {
		if u.HP <= 0 {
			lostUnits++
		}
	}

	moraleThreshold := 0.6 - player.GetTechEffect("morale_threshold")
	if moraleThreshold < 0.1 {
		moraleThreshold = 0.1
	}

	if totalUnits > 0 && float64(lostUnits)/float64(totalUnits) > moraleThreshold {
		for _, u := range nearbyUnits {
			if u.HP > 0 {
				gs.retreatUnit(u)
			}
		}
	}
}

func (gs *GameState) retreatUnit(unit *Unit) {
	gs.removeUnitFromCell(unit)

	player := gs.Players[unit.Owner]
	mainBaseCoord := HexCoord{}
	mainBaseLayer := 0
	found := false

	for key, building := range player.Buildings {
		if building.Type == BuildingMainBase && building.Completed {
			var layerIdx int
			var coordStr string
			fmt.Sscanf(key, "%d_%s", &layerIdx, &coordStr)

			for q := -10; q < 30 && !found; q++ {
				for r := -10; r < 30 && !found; r++ {
					if HexKey(q, r) == coordStr {
						mainBaseCoord = HexCoord{Q: q, R: r}
						mainBaseLayer = layerIdx
						found = true
					}
				}
			}
		}
	}

	if found {
		unit.Layer = mainBaseLayer
		unit.Coord = mainBaseCoord

		layer := gs.Map.Layers[mainBaseLayer]
		if cell, ok := GetCell(layer, mainBaseCoord); ok {
			cell.Units = append(cell.Units, *unit)
		}
	}

	gs.updateUnitInPlayerList(unit)
}
