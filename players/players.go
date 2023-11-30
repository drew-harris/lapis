package players

import (
	"github.com/drew-harris/lapis/graph/model"
	"gorm.io/gorm"
)

func GetAllPlayers(db *gorm.DB) ([]model.Player, error) {

	players := []model.Player{}
	result := db.Find(&players)

	if result.Error != nil {
		return nil, result.Error
	}

	return players, nil
}
