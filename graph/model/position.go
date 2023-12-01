package model

import (
	"time"

	"gorm.io/datatypes"
)

type Position struct {
	Player         *Player        `json:"player" gorm:"not null"`
	PlayerID       string         `json:"playerId" gorm:"not null"`
	ID             string         `json:"id" gorm:"not null"`
	ObjectiveID    string         `json:"objectiveID" gorm:"not null"`
	Unit           string         `json:"unit" gorm:"not null"`
	SavedAt        time.Time      `json:"savedAt" gorm:"not null;autoCreateTime"`
	AdditionalData datatypes.JSON `json:"additionalData,omitempty"`
}
