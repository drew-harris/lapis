package model

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	ID       string  `json:"id"`
	Message  string  `json:"message"`
	Player   *Player `json:"player" gorm:"foreignKey:PlayerID"`
	PlayerID string  `json:"playerId" gorm:"size:255"`
}
