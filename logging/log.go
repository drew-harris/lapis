package logging

import (
	"errors"
	"fmt"
	"time"

	"github.com/drew-harris/lapis/graph/model"
	"github.com/drew-harris/lapis/maps"
	"github.com/google/uuid"
	"github.com/posthog/posthog-go"
	"gorm.io/gorm"
)

type LoggingService struct {
	db      *gorm.DB
	posthog posthog.Client
}

func New(db *gorm.DB, posthog posthog.Client) *LoggingService {
	return &LoggingService{db: db, posthog: posthog}
}

func (l *LoggingService) Log(input model.LogInput) (*model.Log, error) {
	fmt.Println("ABOUT TO LOGG!!!")
	// Check if playerid is valid
	player := model.Player{}
	l.db.Where("id = ?", input.PlayerName).First(&player)
	if l.db.Error != nil {
		fmt.Println(l.db.Error)
		return nil, l.db.Error
	}
	if player.ID == "" {
		return nil, errors.New("Player id is not valid")
	}
	attributes, err := maps.FromMap(input.Attributes)
	if err != nil {
		return nil, err
	}

	// Check if type is empty
	if input.Type == "" {
		return nil, errors.New("Type cannot be empty")
	}

	log := model.Log{
		ID:         uuid.New().String(),
		Message:    input.Message,
		PlayerID:   player.ID,
		Attributes: attributes,
		Type:       input.Type,
		Unit:       input.Unit,
		Objective:  input.Objective,
	}

	fmt.Println("SUPER LOG", log)

	properties := posthog.NewProperties()

	if input.Attributes != nil {
		for key, value := range input.Attributes {
			properties.Set(key, value)
		}
	}

	properties.Set("message", input.Message)

	l.posthog.Enqueue(posthog.Capture{
		DistinctId: log.PlayerID,
		Event:      log.Type.String(),
		Properties: properties,
		Timestamp:  time.Now(),
	})

	fmt.Println("Logging event:", input.Type.String())

	l.db.Create(&log)

	if l.db.Error != nil {
		return nil, l.db.Error
	}

	if err != nil {
		return nil, err
	}

	return &log, nil
}
