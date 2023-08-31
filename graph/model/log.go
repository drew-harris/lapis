package model

import (
	"time"

	"gorm.io/datatypes"
)

type Log struct {
	ID         string         `json:"id" gorm:"primaryKey;size:255"`
	Message    string         `json:"message" gorm:"not null"`
	Player     *Player        `json:"player"`
	PlayerID   string         `json:"playerId" gorm:"size:255" gorm:"not null"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"not null"`
	Attributes datatypes.JSON `json:"attributes"`
	Type       LogType        `json:"string" gorm:"not null"`
}
