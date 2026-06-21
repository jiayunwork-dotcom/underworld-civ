package game

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"
)

type SyncTechsFunc func(gameID string, playerID string, techs []string, currentResearch string, username string) error
type SaveGameRecordFunc func(gameID string, turnNumber int, gameStateJSON string) error
type SaveFinalGameFunc func(gameID string, status string, currentTurn int, winnerID string, endedAt time.Time, players map[string]PlayerStateForSave) error

type PlayerStateForSave struct {
	Score      int
	Eliminated bool
}

type GameManager struct {
	games           map[string]*GameState
	mu              sync.RWMutex
	SyncTechs       SyncTechsFunc
	SaveGameRecord  SaveGameRecordFunc
	SaveFinalGame   SaveFinalGameFunc
}

var manager = &GameManager{
	games: make(map[string]*GameState),
}

func GetGameManager() *GameManager {
	return manager
}

func (gm *GameManager) CreateGame(name string, maxPlayers int, hostID, hostName string, race Race, color string) *GameState {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	seed := time.Now().UnixNano() + rand.Int63()
	state := NewGameState(name, maxPlayers, seed)
	state.AddPlayer(hostID, hostName, race, color, true)

	gm.games[state.ID] = state

	return state
}

func (gm *GameManager) GetGame(gameID string) (*GameState, bool) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	game, ok := gm.games[gameID]
	return game, ok
}

func (gm *GameManager) AddPlayerToGame(gameID, playerID, username string, race Race, color string) (*GameState, bool) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.games[gameID]
	if !ok {
		return nil, false
	}

	if len(game.Players) >= game.MaxPlayers {
		return game, false
	}

	if game.Status != "waiting" {
		return game, false
	}

	game.AddPlayer(playerID, username, race, color, false)
	return game, true
}

func (gm *GameManager) StartGame(gameID string) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.games[gameID]
	if !ok || game.Status != "waiting" {
		return false
	}

	game.StartGame()
	go gm.runGameLoop(gameID)

	return true
}

func (gm *GameManager) runGameLoop(gameID string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	gameEnded := false

	for {
		game, ok := gm.GetGame(gameID)
		if !ok {
			return
		}
		if game.Phase == PhaseEnded && gameEnded {
			return
		}

		<-ticker.C

		gm.mu.Lock()
		if g, ok := gm.games[gameID]; ok {
			if g.Phase == PhasePlanning {
				shouldProcess := false

				if time.Now().After(g.PlanningEndsAt) {
					shouldProcess = true
				} else {
					allReady := true
					hasPlayers := false
					for _, player := range g.Players {
						if !player.Eliminated {
							hasPlayers = true
							if !player.Ready {
								allReady = false
								break
							}
						}
					}
					if hasPlayers && allReady {
						shouldProcess = true
					}
				}

				if shouldProcess {
					g.ProcessTurn()
					gm.syncTechsToRedis(g)
					gm.saveGameRecordToDB(g)

					if g.Phase != PhaseEnded {
						g.PlanningEndsAt = time.Now().Add(90 * time.Second)
						for _, player := range g.Players {
							player.Ready = false
						}
					} else {
						gm.saveFinalGameToDB(g)
						gameEnded = true
					}
				}
			}
		}
		gm.mu.Unlock()
	}
}

func (gm *GameManager) syncTechsToRedis(g *GameState) {
	if gm.SyncTechs == nil {
		return
	}
	for pid, player := range g.Players {
		techList := make([]string, 0)
		for techID, researched := range player.Techs {
			if researched {
				techList = append(techList, techID)
			}
		}
		err := gm.SyncTechs(g.ID, pid, techList, player.CurrentResearch, player.Username)
		if err != nil {
			log.Printf("Warning: Failed to sync techs to Redis for player %s: %v", pid, err)
		}
	}
}

func (gm *GameManager) saveGameRecordToDB(g *GameState) {
	if gm.SaveGameRecord == nil {
		return
	}
	stateJSON, err := json.Marshal(g)
	if err != nil {
		log.Printf("Warning: Failed to marshal game state: %v", err)
		return
	}
	err = gm.SaveGameRecord(g.ID, g.CurrentTurn, string(stateJSON))
	if err != nil {
		log.Printf("Warning: Failed to save game record to DB: %v", err)
	}
}

func (gm *GameManager) saveFinalGameToDB(g *GameState) {
	if gm.SaveFinalGame == nil {
		return
	}

	players := make(map[string]PlayerStateForSave)
	for pid, player := range g.Players {
		players[pid] = PlayerStateForSave{
			Score:      player.Score,
			Eliminated: player.Eliminated,
		}
	}

	err := gm.SaveFinalGame(g.ID, g.Status, g.CurrentTurn, g.WinnerID, time.Now(), players)
	if err != nil {
		log.Printf("Warning: Failed to save final game to DB: %v", err)
	}

	gm.saveGameRecordToDB(g)
}

func (gm *GameManager) SubmitAction(gameID, playerID, action string, data map[string]interface{}) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.games[gameID]
	if !ok || game.Phase != PhasePlanning {
		return false
	}

	player, ok := game.Players[playerID]
	if !ok || player.Eliminated {
		return false
	}

	if action == "ready" {
		player.Ready = true
		return true
	}

	if action == "unready" {
		player.Ready = false
		return true
	}

	if action == "set_research" {
		techID, _ := data["tech_id"].(string)
		ok, _ := game.SetCurrentResearch(playerID, techID)
		return ok
	}

	if action == "blockade_tech" {
		targetID, _ := data["target_id"].(string)
		category, _ := data["category"].(string)
		ok, _ := game.BlockadeTech(playerID, targetID, TechCategory(category))
		return ok
	}

	if action == "dispatch_spy" {
		unitID, _ := data["unit_id"].(string)
		targetID, _ := data["target_id"].(string)
		ok, _ := game.DispatchSpy(playerID, unitID, targetID)
		return ok
	}

	game.SubmitAction(playerID, action, data)
	return true
}

func (gm *GameManager) ListGames() []*GameState {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	games := make([]*GameState, 0, len(gm.games))
	for _, game := range gm.games {
		games = append(games, game)
	}
	return games
}

func (gm *GameManager) RemoveGame(gameID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	delete(gm.games, gameID)
}
