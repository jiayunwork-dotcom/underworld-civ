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

type TechTreeNode struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Cost          int                `json:"cost"`
	Category      game.TechCategory  `json:"category"`
	Tier          int                `json:"tier"`
	Prerequisites []string           `json:"prerequisites"`
	Effects       map[string]float64 `json:"effects"`
	RaceSpecific  game.Race          `json:"race_specific,omitempty"`
	Researched    bool               `json:"researched"`
	Progress      int                `json:"progress"`
	IsCurrent     bool               `json:"is_current"`
	CanResearch   bool               `json:"can_research"`
	PrereqNames   map[string]string  `json:"prereq_names"`
}

func buildTechTree(player *game.PlayerState) []TechTreeNode {
	result := make([]TechTreeNode, 0, len(game.TechDefs))
	for _, tech := range game.TechDefs {
		node := TechTreeNode{
			ID:            tech.ID,
			Name:          tech.Name,
			Description:   tech.Description,
			Cost:          tech.Cost,
			Category:      tech.Category,
			Tier:          tech.Tier,
			Prerequisites: tech.Prerequisites,
			Effects:       tech.Effects,
			RaceSpecific:  tech.RaceSpecific,
			Researched:    false,
			Progress:      0,
			IsCurrent:     false,
			CanResearch:   false,
			PrereqNames:   make(map[string]string),
		}

		if player != nil {
			node.Researched = player.Techs[tech.ID]
			if player.TechProgresses != nil {
				node.Progress = player.TechProgresses[tech.ID]
			}
			node.IsCurrent = player.CurrentResearch == tech.ID
			node.CanResearch = game.CheckPrerequisites(tech.ID, player.Techs) && !node.Researched
		}

		for _, prereqID := range tech.Prerequisites {
			prereqTech := game.GetTechByID(prereqID)
			if prereqTech != nil {
				node.PrereqNames[prereqID] = prereqTech.Name
			}
		}

		result = append(result, node)
	}
	return result
}

type GetTechTreeResponse struct {
	TechDefs         []game.Tech                  `json:"tech_defs"`
	MyTechTree       []TechTreeNode               `json:"my_tech_tree"`
	PlayerResearch   map[string]interface{}       `json:"player_research"`
	OpponentTechs    map[string][]string          `json:"opponent_techs"`
	TechSynergies    map[string]*game.TechSynergy `json:"tech_synergies"`
	KnowledgeReserve int                          `json:"knowledge_reserve"`
	IncomingBlockades []game.TechBlockade         `json:"incoming_blockades"`
	OutgoingBlockades []game.TechBlockade         `json:"outgoing_blockades"`
	OpponentPlayers  []map[string]interface{}     `json:"opponent_players"`
}

func GetTechTree(c *fiber.Ctx) error {
	gameID := c.Params("id")
	playerID := c.Query("player_id")

	state, ok := game.GetGameManager().GetGame(gameID)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "game not found"})
	}

	player := state.Players[playerID]

	opponentTechs := make(map[string][]string)
	playerResearch := make(map[string]interface{})
	opponentPlayers := make([]map[string]interface{}, 0)

	for pid, p := range state.Players {
		techList := make([]string, 0)
		for techID, researched := range p.Techs {
			if researched {
				techList = append(techList, techID)
			}
		}
		opponentTechs[pid] = techList

		playerResearch[pid] = fiber.Map{
			"username":          p.Username,
			"color":             p.Color,
			"race":              p.Race,
			"current_research":  p.CurrentResearch,
			"researched_count":  len(techList),
			"research_points":   p.ResearchPoints,
		}

		if pid != playerID && !p.Eliminated {
			opponentPlayers = append(opponentPlayers, map[string]interface{}{
				"player_id": pid,
				"username":  p.Username,
				"color":     p.Color,
				"race":      p.Race,
			})
		}
	}

	var myTechTree []TechTreeNode
	var synergies map[string]*game.TechSynergy
	var knowledgeReserve int
	var incomingBlockades []game.TechBlockade
	var outgoingBlockades []game.TechBlockade

	if player != nil {
		myTechTree = buildTechTree(player)

		if player.TechSynergies == nil {
			player.TechSynergies = game.CalculateSynergies(player.Techs)
		}

		synergies = make(map[string]*game.TechSynergy)
		for k, v := range player.TechSynergies {
			synergies[string(k)] = v
		}

		knowledgeReserve = player.KnowledgeReserve
		incomingBlockades = player.IncomingBlockades
		outgoingBlockades = player.OutgoingBlockades
	}

	return c.JSON(GetTechTreeResponse{
		TechDefs:          game.TechDefs,
		MyTechTree:        myTechTree,
		PlayerResearch:    playerResearch,
		OpponentTechs:     opponentTechs,
		TechSynergies:     synergies,
		KnowledgeReserve:  knowledgeReserve,
		IncomingBlockades: incomingBlockades,
		OutgoingBlockades: outgoingBlockades,
		OpponentPlayers:   opponentPlayers,
	})
}

type SetResearchRequest struct {
	TechID   string `json:"tech_id"`
	PlayerID string `json:"player_id"`
}

type BlockadeTechRequest struct {
	PlayerID string `json:"player_id"`
	TargetID string `json:"target_id"`
	Category string `json:"category"`
}

func BlockadeTech(c *fiber.Ctx) error {
	gameID := c.Params("id")

	var req BlockadeTechRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.PlayerID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "player id required"})
	}

	if req.TargetID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "target id required"})
	}

	if req.Category == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "category required"})
	}

	state, ok := game.GetGameManager().GetGame(gameID)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "game not found"})
	}

	success, msg := state.BlockadeTech(req.PlayerID, req.TargetID, game.TechCategory(req.Category))
	if !success {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	ws.GetHub().BroadcastToGame(gameID, "game_state", state)

	return c.JSON(fiber.Map{
		"message": "blockade successful",
	})
}

func SetResearch(c *fiber.Ctx) error {
	gameID := c.Params("id")

	var req SetResearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.PlayerID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "player id required"})
	}

	state, ok := game.GetGameManager().GetGame(gameID)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "game not found"})
	}

	success, msg := state.SetCurrentResearch(req.PlayerID, req.TechID)
	if !success {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	ws.GetHub().BroadcastToGame(gameID, "game_state", state)

	player := state.Players[req.PlayerID]
	return c.JSON(fiber.Map{
		"message":          "research target updated",
		"current_research": player.CurrentResearch,
		"research_progress": player.ResearchProgress,
	})
}
