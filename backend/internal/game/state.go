package game

import "time"

type GamePhase string

const (
	PhasePlanning  GamePhase = "planning"
	PhaseExecuting GamePhase = "executing"
	PhaseEnded     GamePhase = "ended"
)

type GameState struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	MaxPlayers     int                      `json:"max_players"`
	CurrentTurn    int                      `json:"current_turn"`
	MaxTurns       int                      `json:"max_turns"`
	Phase          GamePhase                `json:"phase"`
	Status         string                   `json:"status"`
	Map            *GameMap                 `json:"map"`
	Players        map[string]*PlayerState  `json:"players"`
	TurnOrder      []string                 `json:"turn_order"`
	WinnerID       string                   `json:"winner_id"`
	VictoryType    string                   `json:"victory_type"`
	Events         []GameEvent              `json:"events"`
	PendingActions []PlayerAction           `json:"pending_actions"`
	PlanningEndsAt time.Time                `json:"planning_ends_at"`
	LastUpdate     time.Time                `json:"last_update"`
}

type GameEventType string

const (
	EventEarthquake  GameEventType = "earthquake"
	EventFlood       GameEventType = "flood"
	EventMineralVein GameEventType = "mineral_vein"
	EventAncientRuin GameEventType = "ancient_ruin"
	EventRandomGood  GameEventType = "random_good"
)

type GameEvent struct {
	Type     GameEventType `json:"type"`
	Message  string        `json:"message"`
	Turn     int           `json:"turn"`
	Layer    int           `json:"layer"`
	Location HexCoord      `json:"location"`
	Duration int           `json:"duration"`
}

type PlayerAction struct {
	PlayerID string                 `json:"player_id"`
	Action   string                 `json:"action"`
	Data     map[string]interface{} `json:"data"`
}

type MineAction struct {
	Layer int      `json:"layer"`
	Coord HexCoord `json:"coord"`
	Workers int    `json:"workers"`
}

type BuildAction struct {
	Layer    int          `json:"layer"`
	Coord    HexCoord     `json:"coord"`
	Building BuildingType `json:"building"`
}

type MoveAction struct {
	UnitID string   `json:"unit_id"`
	To     HexCoord `json:"to"`
	ToLayer int      `json:"to_layer"`
}

type AttackAction struct {
	UnitID   string   `json:"unit_id"`
	Target   HexCoord `json:"target"`
	TargetLayer int  `json:"target_layer"`
}
