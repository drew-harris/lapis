package code

import (
	"math/rand"

	"github.com/drew-harris/lapis/graph/model"
	"gorm.io/gorm"
)

// Function to get six random letters
func getRandomLetters() string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RegisterPlayerWithNewCode(name string, db *gorm.DB) (*model.Player, error) {
	player := model.Player{}

	var id string
	// TODO: Replace with mojang api call
	id = getRandomLetters()

	// Make sure code is unique
	for {
		var count int64
		db.Model(&model.Player{}).Where("id = ?", id).Count(&count)
		if count == 0 {
			break
		}
		id = getRandomLetters()
	}

	player = model.Player{
		ID:   id,
		Name: name,
	}

	result := db.Create(&player)
	if result.Error != nil {
		return nil, result.Error
	}
	return &player, nil
}

type CreatePlayerInput struct {
	Name string `json:"name"`
}
