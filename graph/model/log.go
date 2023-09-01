package model

import (
	"time"

	"gorm.io/datatypes"
)

type Log struct {
	ID         string         `json:"id" gorm:"primaryKey;size:255"`
	Message    *string        `json:"message"` // Can be null
	Player     *Player        `json:"player" gorm:"not null"`
	PlayerID   string         `json:"playerId" gorm:"size:255;not null"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"not null"`
	Attributes datatypes.JSON `json:"attributes" gorm:"not null"`
	Type       LogType        `json:"string" gorm:"not null"`
}
