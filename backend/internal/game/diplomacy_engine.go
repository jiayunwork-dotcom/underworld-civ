package game

import (
	"github.com/google/uuid"
)

func (gs *GameState) processDiplomacy() {
	for _, action := range gs.PendingActions {
		switch action.Action {
		case "trade_offer":
			gs.processTradeOffer(action)
		case "accept_trade":
			gs.processAcceptTrade(action)
		case "propose_alliance":
			gs.processProposeAlliance(action)
		case "declare_war":
			gs.processDeclareWar(action)
		case "propose_passage":
			gs.processProposePassage(action)
		}
	}

	for _, player := range gs.Players {
		if player.Eliminated {
			continue
		}
		for _, rel := range player.Diplomacy {
			if rel.CooldownTurns > 0 {
				rel.CooldownTurns--
			}
			if rel.TreatyViolationPenalty > 0 {
				rel.TreatyViolationPenalty--
			}
		}
	}
}

func (gs *GameState) processTradeOffer(action PlayerAction) {
	fromPlayer := action.PlayerID
	toPlayer := action.Data["to_player"].(string)

	offerData := action.Data["offer"].(map[string]interface{})
	demandData := action.Data["demand"].(map[string]interface{})

	offer := Resources{
		Stone:        int(offerData["stone"].(float64)),
		Metal:        int(offerData["metal"].(float64)),
		GlowMushroom: int(offerData["glow_mushroom"].(float64)),
		Water:        int(offerData["water"].(float64)),
		MagicCrystal: int(offerData["magic_crystal"].(float64)),
		FossilFuel:   int(offerData["fossil_fuel"].(float64)),
	}
	demand := Resources{
		Stone:        int(demandData["stone"].(float64)),
		Metal:        int(demandData["metal"].(float64)),
		GlowMushroom: int(demandData["glow_mushroom"].(float64)),
		Water:        int(demandData["water"].(float64)),
		MagicCrystal: int(demandData["magic_crystal"].(float64)),
		FossilFuel:   int(demandData["fossil_fuel"].(float64)),
	}

	from := gs.Players[fromPlayer]
	to := gs.Players[toPlayer]

	if from == nil || to == nil || from.Eliminated || to.Eliminated {
		return
	}

	if !gs.canAfford(from, offer) {
		return
	}

	tradeOffer := TradeOffer{
		ID:         uuid.New().String(),
		FromPlayer: fromPlayer,
		ToPlayer:   toPlayer,
		Offer:      offer,
		Demand:     demand,
		TurnsValid: 3,
	}

	to.TradeOffers = append(to.TradeOffers, tradeOffer)
}

func (gs *GameState) processAcceptTrade(action PlayerAction) {
	player := gs.Players[action.PlayerID]
	if player == nil || player.Eliminated {
		return
	}

	tradeID := action.Data["trade_id"].(string)

	var trade *TradeOffer
	for i := range player.TradeOffers {
		if player.TradeOffers[i].ID == tradeID {
			trade = &player.TradeOffers[i]
			break
		}
	}

	if trade == nil {
		return
	}

	fromPlayer := gs.Players[trade.FromPlayer]
	if fromPlayer == nil || fromPlayer.Eliminated {
		return
	}

	if !gs.canAfford(fromPlayer, trade.Offer) {
		return
	}
	if !gs.canAfford(player, trade.Demand) {
		return
	}

	gs.payCost(fromPlayer, trade.Offer)
	gs.payCost(player, trade.Demand)

	player.Resources.Stone += trade.Offer.Stone
	player.Resources.Metal += trade.Offer.Metal
	player.Resources.GlowMushroom += trade.Offer.GlowMushroom
	player.Resources.Water += trade.Offer.Water
	player.Resources.MagicCrystal += trade.Offer.MagicCrystal
	player.Resources.FossilFuel += trade.Offer.FossilFuel

	fromPlayer.Resources.Stone += trade.Demand.Stone
	fromPlayer.Resources.Metal += trade.Demand.Metal
	fromPlayer.Resources.GlowMushroom += trade.Demand.GlowMushroom
	fromPlayer.Resources.Water += trade.Demand.Water
	fromPlayer.Resources.MagicCrystal += trade.Demand.MagicCrystal
	fromPlayer.Resources.FossilFuel += trade.Demand.FossilFuel

	newOffers := make([]TradeOffer, 0)
	for _, o := range player.TradeOffers {
		if o.ID != tradeID {
			newOffers = append(newOffers, o)
		}
	}
	player.TradeOffers = newOffers
}

func (gs *GameState) processProposeAlliance(action PlayerAction) {
	fromPlayer := action.PlayerID
	toPlayer := action.Data["to_player"].(string)

	from := gs.Players[fromPlayer]
	to := gs.Players[toPlayer]

	if from == nil || to == nil || from.Eliminated || to.Eliminated {
		return
	}

	rel := gs.getOrCreateRelation(fromPlayer, toPlayer)
	rel.Status = DiplomacyAlliance
	rel.TradeAccess = true
	rel.PassageRights = true
	rel.SharedVision = true

	rel2 := gs.getOrCreateRelation(toPlayer, fromPlayer)
	rel2.Status = DiplomacyAlliance
	rel2.TradeAccess = true
	rel2.PassageRights = true
	rel2.SharedVision = true
}

func (gs *GameState) processDeclareWar(action PlayerAction) {
	fromPlayer := action.PlayerID
	toPlayer := action.Data["to_player"].(string)

	from := gs.Players[fromPlayer]
	to := gs.Players[toPlayer]

	if from == nil || to == nil || from.Eliminated || to.Eliminated {
		return
	}

	rel := gs.getOrCreateRelation(fromPlayer, toPlayer)
	if rel.Status == DiplomacyAlliance {
		rel.TreatyViolationPenalty = 5
	}
	rel.Status = DiplomacyWar
	rel.CooldownTurns = 3
	rel.TradeAccess = false
	rel.PassageRights = false
	rel.SharedVision = false

	rel2 := gs.getOrCreateRelation(toPlayer, fromPlayer)
	if rel2.Status == DiplomacyAlliance {
		rel2.TreatyViolationPenalty = 5
	}
	rel2.Status = DiplomacyWar
	rel2.CooldownTurns = 3
	rel2.TradeAccess = false
	rel2.PassageRights = false
	rel2.SharedVision = false
}

func (gs *GameState) processProposePassage(action PlayerAction) {
	fromPlayer := action.PlayerID
	toPlayer := action.Data["to_player"].(string)

	from := gs.Players[fromPlayer]
	to := gs.Players[toPlayer]

	if from == nil || to == nil || from.Eliminated || to.Eliminated {
		return
	}

	rel := gs.getOrCreateRelation(fromPlayer, toPlayer)
	rel.PassageRights = true
	rel.TradeAccess = true

	rel2 := gs.getOrCreateRelation(toPlayer, fromPlayer)
	rel2.PassageRights = true
	rel2.TradeAccess = true
}

func (gs *GameState) getOrCreateRelation(playerA, playerB string) *DiplomaticRelation {
	player := gs.Players[playerA]
	if player == nil {
		return nil
	}

	if rel, ok := player.Diplomacy[playerB]; ok {
		return rel
	}

	rel := &DiplomaticRelation{
		PlayerA: playerA,
		PlayerB: playerB,
		Status:  DiplomacyNeutral,
	}
	player.Diplomacy[playerB] = rel
	return rel
}

func (gs *GameState) SubmitAction(playerID, action string, data map[string]interface{}) {
	gs.PendingActions = append(gs.PendingActions, PlayerAction{
		PlayerID: playerID,
		Action:   action,
		Data:     data,
	})
}
