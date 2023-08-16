package model

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	ID   string `json:"id" gorm:"primaryKey" gorm:"size:255"`
	Name string `json:"name"`
	Logs []*Log `json:"logs"`
}
