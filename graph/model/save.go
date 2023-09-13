package model

import (
	"time"

	"gorm.io/datatypes"
)

type Save struct {
	ID        string         `json:"id" gorm:"primaryKey;size:255"`
	Name      string         `json:"name" gorm:"not null"`
	Player    *Player        `json:"player" gorm:"not null"`
	PlayerID  string         `json:"playerId" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time      `json:"lastSavedAt" gorm:"not null"`
	GraphData datatypes.JSON `json:"graphData"`
}
