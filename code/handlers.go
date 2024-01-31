package code

import (
	"time"

	"github.com/drew-harris/lapis/players"
	"github.com/drew-harris/lapis/views"
	"github.com/gofiber/fiber/v2"
	"github.com/posthog/posthog-go"
	"gorm.io/gorm"
)

type CodeHandler struct {
	db      *gorm.DB
	posthog posthog.Client
}

func CreateCodeHandler(db *gorm.DB, posthog posthog.Client) *CodeHandler {
	return &CodeHandler{
		db:      db,
		posthog: posthog,
	}
}

// GET /codes
func (h *CodeHandler) GetCodes(c *fiber.Ctx) error {
	players, err := players.GetAllPlayers(h.db)
	if err != nil {
		return err
	}
	return views.SendView(c, views.CodePage(players))
}

// POST hx/setup
func (h *CodeHandler) SetupPlayer(c *fiber.Ctx) error {
	playerInput := CreatePlayerInput{}
	if err := c.BodyParser(&playerInput); err != nil {
		return err
	}

	player, err := RegisterPlayerWithNewCode(playerInput, h.db)
	if err != nil {
		return err
	}

	h.posthog.Enqueue(posthog.Alias{
		DistinctId: player.ID,
		Alias:      player.Name,
		Timestamp:  time.Now(),
	})

	h.posthog.Enqueue(posthog.Identify{
		DistinctId: player.ID,
		Properties: posthog.NewProperties().Set("name", player.Name).Set("created", time.Now()),
		Timestamp:  time.Now(),
	})

	return views.SendView(c, views.PlayerRow(*player))
}
