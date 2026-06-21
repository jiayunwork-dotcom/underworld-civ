package handlers

import (
	"net/http"
	"underworld-civ/internal/game"
	"underworld-civ/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateGameRequest struct {
	Name       string `json:"name"`
	MaxPlayers int    `json:"max_players"`
	Race       string `json:"race"`
	Color      string `json:"color"`
	PlayerID   string `json:"player_id"`
	Username   string `json:"username"`
}

type JoinGameRequest struct {
	Race     string `json:"race"`
	Color    string `json:"color"`
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
}

func CreateGame(c *fiber.Ctx) error {
	var req CreateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.MaxPlayers < 4 || req.MaxPlayers > 8 {
		req.MaxPlayers = 6
	}

	playerID := req.PlayerID
	username := req.Username

	if playerID == "" {
		playerID = uuid.New().String()
	}
	if username == "" {
		username = "Player_" + playerID[:8]
	}

	race := game.Race(req.Race)
	if _, ok := game.RaceDefs[race]; !ok {
		race = game.RaceDwarf
	}

	if req.Color == "" {
		req.Color = "#e74c3c"
	}

	state := game.GetGameManager().CreateGame(req.Name, req.MaxPlayers, playerID, username, race, req.Color)
	ws.GetHub().CreateGameRoom(state.ID, state)

	return c.JSON(fiber.Map{
		"game_id":   state.ID,
		"player_id": playerID,
		"username":  username,
		"game":      state,
	})
}

func ListGames(c *fiber.Ctx) error {
	games := game.GetGameManager().ListGames()

	result := make([]fiber.Map, 0, len(games))
	for _, g := range games {
		players := make([]fiber.Map, 0)
		for _, p := range g.Players {
			players = append(players, fiber.Map{
				"id":       p.PlayerID,
				"username": p.Username,
				"race":     p.Race,
				"color":    p.Color,
				"is_host":  p.IsHost,
			})
		}
		result = append(result, fiber.Map{
			"id":           g.ID,
			"name":         g.Name,
			"status":       g.Status,
			"max_players":  g.MaxPlayers,
			"player_count": len(g.Players),
			"current_turn": g.CurrentTurn,
			"players":      players,
		})
	}

	return c.JSON(result)
}

func GetGame(c *fiber.Ctx) error {
	gameID := c.Params("id")

	state, ok := game.GetGameManager().GetGame(gameID)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "game not found"})
	}

	playerID := c.Query("player_id")

	return c.JSON(fiber.Map{
		"game":      state,
		"player_id": playerID,
	})
}

func JoinGame(c *fiber.Ctx) error {
	gameID := c.Params("id")

	var req JoinGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	playerID := req.PlayerID
	username := req.Username

	if playerID == "" {
		playerID = uuid.New().String()
	}
	if username == "" {
		username = "Player_" + playerID[:8]
	}

	race := game.Race(req.Race)
	if _, ok := game.RaceDefs[race]; !ok {
		race = game.RaceDwarf
	}

	if req.Color == "" {
		req.Color = "#3498db"
	}

	state, success := game.GetGameManager().AddPlayerToGame(gameID, playerID, username, race, req.Color)
	if !success {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "could not join game"})
	}

	hub := ws.GetHub()
	if _, ok := hub.GetGameRoom(gameID); !ok {
		hub.CreateGameRoom(gameID, state)
	}

	return c.JSON(fiber.Map{
		"game_id":   gameID,
		"player_id": playerID,
		"username":  username,
		"game":      state,
	})
}

func StartGame(c *fiber.Ctx) error {
	gameID := c.Params("id")

	success := game.GetGameManager().StartGame(gameID)
	if !success {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "could not start game"})
	}

	state, _ := game.GetGameManager().GetGame(gameID)

	hub := ws.GetHub()
	if _, ok := hub.GetGameRoom(gameID); !ok {
		hub.CreateGameRoom(gameID, state)
	}
	hub.BroadcastToGame(gameID, "game_state", state)

	return c.JSON(fiber.Map{
		"message": "game started",
		"game":    state,
	})
}

func SubmitAction(c *fiber.Ctx) error {
	gameID := c.Params("id")

	var action struct {
		Action   string                 `json:"action"`
		Data     map[string]interface{} `json:"data"`
		PlayerID string                 `json:"player_id"`
	}

	if err := c.BodyParser(&action); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	playerID := action.PlayerID
	if playerID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "player id required"})
	}

	success := game.GetGameManager().SubmitAction(gameID, playerID, action.Action, action.Data)
	if !success {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "could not submit action"})
	}

	state, _ := game.GetGameManager().GetGame(gameID)
	ws.GetHub().BroadcastToGame(gameID, "game_state", state)

	return c.JSON(fiber.Map{"message": "action submitted"})
}

func GetRaces(c *fiber.Ctx) error {
	races := make([]game.RaceInfo, 0, len(game.RaceDefs))
	for _, race := range game.RaceDefs {
		races = append(races, race)
	}
	return c.JSON(races)
}

func GetBuildings(c *fiber.Ctx) error {
	buildings := make([]game.BuildingInfo, 0, len(game.BuildingDefs))
	for _, building := range game.BuildingDefs {
		buildings = append(buildings, building)
	}
	return c.JSON(buildings)
}

func GetUnits(c *fiber.Ctx) error {
	units := make([]game.UnitInfo, 0, len(game.UnitDefs))
	for _, unit := range game.UnitDefs {
		units = append(units, unit)
	}
	return c.JSON(units)
}

func GetTechs(c *fiber.Ctx) error {
	return c.JSON(game.TechDefs)
}
