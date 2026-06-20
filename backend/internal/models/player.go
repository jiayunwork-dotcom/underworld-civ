package models

import "time"

type Player struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	Username  string    `json:"username" gorm:"unique;size:50"`
	Password  string    `json:"-" gorm:"column:password_hash;size:255"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}

func (Player) TableName() string {
	return "players"
}
