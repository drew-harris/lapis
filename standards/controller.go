package standards

import (
	"fmt"

	"github.com/drew-harris/lapis/graph/model"
	"github.com/drew-harris/lapis/logging"
	"github.com/drew-harris/lapis/realtime"
	"github.com/drew-harris/lapis/views"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type StandardsController struct {
	Db       *gorm.DB
	Realtime *realtime.Service
	Logging  *logging.LoggingService
}

func (s *StandardsController) ShowStandardsPage(c *fiber.Ctx) error {
	return views.SendView(c, views.StandardsPage())
}

func (s *StandardsController) SendTestLog(c *fiber.Ctx) error {
	fakeNumber := "3"

	log, err := s.Logging.Log(model.LogInput{
		Type:       model.LogTypeTestLog,
		Unit:       &fakeNumber,
		Objective:  &fakeNumber,
		PlayerName: "DREW",
		Attributes: map[string]interface{}{},
		Message:    nil,
	})

	if err != nil {
		views.SendView(c, views.Error("Could not create test message"))
	}

	fmt.Println("Sending log through realtime")

	err = s.Realtime.SendHtml(views.LogRow(*log))
	if err != nil {
		fmt.Println("realtime error")
		return nil
	}

	return nil
}
