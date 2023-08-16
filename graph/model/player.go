package model

import "time"

type Player struct {
	ID        string    `json:"id" gorm:"primaryKey;not null;size:255"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	Logs      []Log     `json:"logs" gorm:"foreignKey:PlayerID;references:ID"`
}
