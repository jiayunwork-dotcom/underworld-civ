package models

import "time"

type Game struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string    `json:"name" gorm:"size:100"`
	MaxPlayers  int       `json:"max_players" gorm:"default:6"`
	Status      string    `json:"status" gorm:"size:20;default:'waiting'"`
	CurrentTurn int       `json:"current_turn" gorm:"default:0"`
	MaxTurns    int       `json:"max_turns" gorm:"default:50"`
	WinnerID    string    `json:"winner_id" gorm:"type:uuid"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
}

func (Game) TableName() string {
	return "games"
}

type GamePlayer struct {
	ID         string `json:"id" gorm:"primaryKey;type:uuid"`
	GameID     string `json:"game_id" gorm:"type:uuid"`
	PlayerID   string `json:"player_id" gorm:"type:uuid"`
	Race       string `json:"race" gorm:"size:30"`
	Color      string `json:"color" gorm:"size:20"`
	IsHost     bool   `json:"is_host" gorm:"default:false"`
	Score      int    `json:"score" gorm:"default:0"`
	Eliminated bool   `json:"eliminated" gorm:"default:false"`
	Username   string `json:"username" gorm:"-"`
}

func (GamePlayer) TableName() string {
	return "game_players"
}
