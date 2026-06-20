package game

import (
	"math/rand"
	"sync"
	"time"
)

type GameManager struct {
	games map[string]*GameState
	mu    sync.RWMutex
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
	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		game, ok := gm.GetGame(gameID)
		if !ok || game.Phase == PhaseEnded {
			return
		}

		<-ticker.C

		gm.mu.Lock()
		if g, ok := gm.games[gameID]; ok {
			if g.Phase == PhasePlanning {
				g.ProcessTurn()

				if g.Phase != PhaseEnded {
					g.PlanningEndsAt = time.Now().Add(90 * time.Second)
				}
			}
		}
		gm.mu.Unlock()
	}
}

func (gm *GameManager) SubmitAction(gameID, playerID, action string, data map[string]interface{}) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.games[gameID]
	if !ok || game.Phase != PhasePlanning {
		return false
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
