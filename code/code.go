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

func RegisterPlayerWithNewCode(input CreatePlayerInput, db *gorm.DB) (*model.Player, error) {
	player := model.Player{}

	if input.Code == "" {
		input.Code = getRandomLetters()
	}

	// Make sure code is unique
	for {
		var count int64
		db.Model(&model.Player{}).Where("id = ?", input.Code).Count(&count)
		if count == 0 {
			break
		}
		input.Code = getRandomLetters()
	}

	player = model.Player{
		ID:   input.Code,
		Name: input.Name,
	}

	result := db.Create(&player)
	if result.Error != nil {
		return nil, result.Error
	}
	return &player, nil
}

type CreatePlayerInput struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
